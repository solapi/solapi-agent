package main

import (
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"
  "unsafe"
  "encoding/json"
  "io/ioutil"
  "database/sql"

  "github.com/takama/daemon"
  _ "github.com/go-sql-driver/mysql"
  "github.com/solapi/solapi-go"
)

type DBConfig struct {
  Provider string `json:"provider"`
  DBName string `json:"dbname"`
  Table string `json:"table"`
  User string `json:"user"`
  Password string `json:"password"`
  Host string `json:"host"`
  Port int `json:"port"`
}

const (
  name        = "solapi"
  description = "Solapi Agent Service"
)

var dbconf DBConfig

var stdlog, errlog *log.Logger

var client *solapi.Client

var db *sql.DB

// Service has embedded daemon
type Service struct {
  daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

  usage := "Usage: solapi install | remove | start | stop | status"

  // if received any kind of command, do it
  if len(os.Args) > 1 {
    command := os.Args[1]
    switch command {
    case "install":
      return service.Install()
    case "remove":
      return service.Remove()
    case "start":
      return service.Start()
    case "stop":
      return service.Stop()
    case "status":
      return service.Status()
    default:
      return usage, nil
    }
  }

  // Set up channel on which to send signal notifications.
  // We must use a buffered channel or risk missing the signal
  // if we're not ready to receive when the signal is sent.
  interrupt := make(chan os.Signal, 1)
  signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

  // loop work cycle with accept connections or interrupt
  // by system signal

  var err error

  var b []byte
	b, err = ioutil.ReadFile("./db.json")
	if err != nil {
		fmt.Println(err)
		return "db.json 로딩 오류", nil
	}
	json.Unmarshal(b, &dbconf)
  connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.DBName)
  fmt.Println(connectionString)

  client = solapi.NewClient()
  db, err = sql.Open("mysql", connectionString)
  if err != nil {
    panic(err)
  }
  db.SetConnMaxLifetime(time.Minute * 3)
  db.SetMaxOpenConns(10)
  db.SetMaxIdleConns(10)

  go pollMsg()
  go pollResult()
  go pollLastReport()

  for {
    select {
    case killSignal := <-interrupt:
      stdlog.Println("Got signal:", killSignal)
      if killSignal == os.Interrupt {
        return "Daemon was interrupted by system signal", nil
      }
      return "Daemon was killed", nil
    }
  }

  return usage, nil
}

func pollMsg() {
  for {
    stdlog.Println("Polling Msg...")
    time.Sleep(time.Second * 1)
    rows, err := db.Query("SELECT id, payload FROM msg WHERE sent = false AND sendAttempts < 3")
    if err != nil {
      fmt.Println(err)
    }
    defer rows.Close()
    if err != nil {
      fmt.Println(err)
    }

    var messageList []interface{}
    var groupId string 

    var id uint32
    var payload string
    var msgObj interface{}
    var count int = 0
    var idList []uint32
    for rows.Next() {
      if count == 0 {
        params := make(map[string]string)
        result, err := client.Messages.CreateGroup(params)
        printObj(result)
        fmt.Println(result.GroupId)
        if err != nil {
          fmt.Println(err)
        }
        groupId = result.GroupId
      }
      count++

      err := rows.Scan(&id, &payload)
      if err != nil {
        fmt.Println(err)
      }
      fmt.Sprintf("id: %u", id)
      _, err = db.Exec("UPDATE msg SET sent = true WHERE id = ?", id)
      if err != nil {
        fmt.Println(err)
        continue
      }
      err = json.Unmarshal([]byte(payload), &msgObj)
      if err != nil {
        fmt.Println(err)
        continue
      }

      messageList = append(messageList, msgObj)
      idList = append(idList, id)
    }
    if len(messageList) > 0 {
      var msgParams = make(map[string]interface{})
      msgParams["messages"] = messageList

      result2, err := client.Messages.AddGroupMessage(groupId, msgParams)
      if err != nil {
        fmt.Println(err)
        continue
      }
      printObj(result2)
      for i, res := range(result2.ResultList) {
        _, err = db.Exec("UPDATE msg SET result = json_object('messageId', ?, 'groupId', ?, 'statusCode', ?, 'statusMessage', ?) WHERE id = ?", res.MessageId, groupId, res.StatusCode, res.StatusMessage,idList[i])
        if err != nil {
          fmt.Println(err)
          continue
        }
      }

      _, err = client.Messages.SendGroup(groupId)
      if err != nil {
        fmt.Println(err)
        continue
      }
    }
  }
}

func pollLastReport() {
  for {
    stdlog.Println("Polling Last Report...")
    time.Sleep(time.Second * 2)
    rows, err := db.Query("SELECT id, messageId, statusCode FROM msg WHERE sent = true AND createdAt > SUBDATE(NOW(), INTERVAL 72 HOUR) AND statusCode IN ('2000', '3000')")
    if err != nil {
      fmt.Println(err)
    }
    defer rows.Close()

    var id uint32
    var messageId string
    var statusCode string
    for rows.Next() {
      rows.Scan(&id, &messageId, &statusCode)
      if syncMsgStatus(db, messageId, statusCode) == false {
        _, err = db.Exec("UPDATE msg SET result = json_set(result, '$.statusCode', ?), updatedAt = NOW() WHERE id = ?", "3040", id)
      }
    }
  }
}

func pollResult() {
  for {
    stdlog.Println("Polling Result...")
    time.Sleep(time.Second * 2)
    rows, err := db.Query("SELECT id, messageId, statusCode FROM msg WHERE sent = true AND createdAt > SUBDATE(NOW(), INTERVAL 72 HOUR) AND updatedAt < SUBDATE(NOW(), INTERVAL (10 * (reportAttempts + 1)) SECOND) AND reportAttempts < 10 AND statusCode IN ('2000', '3000')")
    if err != nil {
      fmt.Println(err)
    }
    defer rows.Close()

    var id uint32
    var messageId string
    var statusCode string
    for rows.Next() {
      rows.Scan(&id, &messageId, &statusCode)

      _, err = db.Exec("UPDATE msg SET reportAttempts = reportAttempts + 1, updatedAt = NOW() WHERE id = ?", id)

      syncMsgStatus(db, messageId, statusCode)
    }
  }
}

func syncMsgStatus(db *sql.DB, messageId string, statusCode string) bool {
  params := make(map[string]string)
  params["criteria"] = "messageId"
  params["cond"] = "eq"
  params["value"] = messageId

  result, err := client.Messages.GetMessageList(params)
  printObj(result)
  if err != nil {
    fmt.Println(err)
  }

  for i, res := range(result.MessageList) {
    fmt.Println(i)
    printObj(res)
    if res.StatusCode != statusCode {
      _, err = db.Exec("UPDATE msg SET result = json_set(result, '$.statusCode', ?), updatedAt = NOW() WHERE messageId = ?", res.StatusCode, res.MessageId)
      if err != nil {
        panic(err)
      }
      return true
    }
  }

  return false
}

func printObj(obj interface{}) {
  var msgBytes []byte
  msgBytes, err := json.MarshalIndent(obj, "", "  ")
  if err != nil {
    panic(err)
  }
  msgStr := *(*string)(unsafe.Pointer(&msgBytes))
  fmt.Println(msgStr)
}

func init() {
  stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
  errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)

}

func main() {
  srv, err := daemon.New(name, description, daemon.UserAgent)
  if err != nil {
    fmt.Println("test")
    errlog.Println("Error: ", err)
    os.Exit(1)
  }
  service := &Service{srv}
  status, err := service.Manage()
  if err != nil {
    errlog.Println(status, "\nError: ", err)
    os.Exit(1)
  }
  fmt.Println(status)
}
