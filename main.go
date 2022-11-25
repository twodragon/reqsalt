package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"strings"

	"os"

	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"gopkg.in/guregu/null.v3"
)


var (
	ip          = "127.0.0.1"
	port        = "5477"
	user        = "postgres"
	pass        = ""
	dbname      = "beta_1ch"
	//dinle 		= "1010478927162114078"
	BotID 	 string
)

func main(){
StartLogging()
Connect()
}

func StartLogging() {
	fi, err := os.OpenFile("RegLog.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //log file
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(fi)
}

func Connect() {
	token := "   "

	// Create a new Discord session using the provided bot token.
	goBot, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
	}
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	} 
	BotID = u.ID
	goBot.AddHandler(isTheChannelTheMessageWasSentInPrivate)
	// Open a websocket connection to Discord and begin listening.
	err = goBot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	goBot.Close()
}
func isTheChannelTheMessageWasSentInPrivate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	}

 if strings.Contains(m.Content ,"!register") {
//	if m.ChannelID == dinle {

		parts := strings.Split(m.Content, " ");
		
		if len(parts) < 4 {
			s.ChannelMessageSend(m.ChannelID, "hata oluştu eksik yazdın\n kayıt olmak için mesajını şu şekilde yaz \n !register id password mail    ")
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			return 
		}

		userid := parts[1];
		passw := parts[2];
		mail := parts[3]
		role := "1"
	    ctime := null.NewTime(time.Now(), true)
		data := []byte(passw)
		hashx := sha256.Sum256(data)
		hasp :=  string(hashx[:])
		salt 		:= "???"
        test := fmt.Sprintf("%X%s", hasp, salt) //salt test --
		data2 := []byte(test)
		hashx2 := sha256.Sum256(data2)
		hasp2 := string(hashx2[:])
		saltpw := fmt.Sprintf("%X", hasp2) //--
		psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ip, port, user, pass, dbname)

	
		postgres_insert_query := `INSERT INTO hops.users (user_name, password, mail, user_type, created_at) VALUES ($1,$2,$3,$4,$5)`
                
				   db, err := sql.Open("postgres", psqlconn)
				   CheckError(err)
				   err = db.Ping()
                   CheckError(err)
				   _, err = db.Exec(postgres_insert_query, userid, saltpw, mail, role, ctime)
				   if err != nil {
				
					s.ChannelMessageSend(m.ChannelID, "hata oluştu sanırım bu id başkası tarafından kullanılıyor")
					s.ChannelMessageDelete(m.ChannelID, m.ID)
					
					log.Printf( "ERROR %s şifre %s mail %s", userid, passw, mail )
					//CheckError(err)
					return
				}else{
					s.ChannelMessageSend(m.ChannelID, "success!\n id: "+ fmt.Sprintf( "%s \nşifre: %s \nmail %s", userid, passw, mail)  )
					s.ChannelMessageDelete(m.ChannelID, m.ID)
					log.Printf( "%s şifre: %s mail: %s", userid, passw, mail )
				}
		s.ChannelMessageDelete(m.ChannelID, m.ID)

	}
	//}else{
	//	return
	//}
}

func CheckError(err error) {
	if err != nil {
	log.Fatal(err)
	}
}

func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}