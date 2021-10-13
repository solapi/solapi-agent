package main

import (
  "fmt"
  "os"
  "time"
  "unsafe"
  "strconv"
  "encoding/json"
  "io/ioutil"
  "database/sql"

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

type APIConfig struct {
  APIKey     string `json:"apiKey"`
  APISecret  string `json:"APISecret"`
  Protocol   string `json:"Protocol"`
  Domain     string `json:"Domain"`
  Prefix     string `json:"Prefix"`
  AppId      string `json:"AppId"`
}

var dbconf DBConfig
var apiconf APIConfig

var client *solapi.Client

var db *sql.DB

var homedir string = "/opt/agent"

func getConnectionString(homedir string) (string, error) {
  var b []byte
  fmt.Println(homedir)
	b, err := ioutil.ReadFile(homedir + "/db.json")
	if err != nil {
		fmt.Println(err)
		return "db.json 로딩 오류", err
	}
	json.Unmarshal(b, &dbconf)
  connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.DBName)
  fmt.Println(connectionString)
  return connectionString, nil
}

func getAPIConfig(homedir string, apiconf *APIConfig) error {
  var b []byte
	b, err := ioutil.ReadFile(homedir + "/config.json")
	if err != nil {
		fmt.Println(err)
    return err
	}
	json.Unmarshal(b, &apiconf)
  return nil
}


func syncMsgStatus(messageIds []string) {
  b, _ := json.Marshal(messageIds)
  params := make(map[string]string)
  params["messageIds[in]"] = string(b)
  params["limit"] = strconv.Itoa(len(messageIds))

  fmt.Println("메시지 상태 동기화:", len(messageIds), "건")

  result, err := client.Messages.GetMessageList(params)
  if err != nil {
    fmt.Println(err)
  }

  for _, res := range(result.MessageList) {
    _, err = db.Exec("UPDATE msg SET result = json_set(result, '$.status', ?, '$.statusCode', ?, '$.statusMessage', ?), updatedAt = NOW() WHERE messageId = ?", res.Status, res.StatusCode, res.Reason, res.MessageId)
    if err != nil {
      panic(err)
    }
  }
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

func main() {
  agentHome := os.Getenv("AGENT_HOME")
  if len(agentHome) > 0 {
    homedir = agentHome
  }

  var err error

  connectionString, _ := getConnectionString(homedir)
  fmt.Println("DB에 연결합니다 Connection String:", connectionString)

  db, err = sql.Open("mysql", connectionString)
  if err != nil {
    panic(err)
  }
  db.SetConnMaxLifetime(time.Minute * 3)
  db.SetMaxOpenConns(10)
  db.SetMaxIdleConns(10)

  err = getAPIConfig(homedir, &apiconf)
  if err != nil {
    panic(err)
  }

  client = solapi.NewClient()
	client.Messages.Config = map[string]string{
	  "APIKey": apiconf.APIKey,
    "APISecret": apiconf.APISecret,
    "Protocol": apiconf.Protocol,
    "Domain": apiconf.Domain,
    "Prefix": apiconf.Prefix,
	}
	client.Storage.Config = map[string]string{
	  "APIKey": apiconf.APIKey,
    "APISecret": apiconf.APISecret,
    "Protocol": apiconf.Protocol,
    "Domain": apiconf.Domain,
    "Prefix": apiconf.Prefix,
	}


  rows, err := db.Query("SELECT id, messageId FROM msg WHERE sent = true AND status = 'PENDING'")
  if err != nil {
    fmt.Println("DB Query ERROR:", err)
    return
  }
  defer rows.Close()

  var id uint32
  var messageId string
  var messageIds []string
  for rows.Next() {
    rows.Scan(&id, &messageId)
    messageIds = append(messageIds, messageId)
  }
  if len(messageIds) > 0 {
    syncMsgStatus(messageIds)
  }
}
