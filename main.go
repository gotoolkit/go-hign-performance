package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/llitfkitfk/GoHighPerformance/pkg/db"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/redis.v5"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
	"github.com/llitfkitfk/GoHighPerformance/pkg/handler"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type User struct {
	Id string
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Authenticate the request, get the id from the route params,
	// and fetch the user from the DB, etc.

	log.Print("Get User")
	// Get the token and pass it in the CSRF header. Our JSON-speaking client
	// or JavaScript framework can now read the header and return the token in
	// in its own "X-CSRF-Token" request header on the subsequent POST.
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	var user User
	b, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(b)
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Print("Home")
	fmt.Fprintf(w, "Hello World, %s!", "test")
}

func ginDemo() {
	gin.Default()
}

func main() {
	// Our top-level router doesn't need CSRF protection: it's simple.

	// ... but our /api/* routes do, so we add it to the sub-router only.

	//r.Use(csrf.Protect([]byte("32-byte-long-auth-key")))

	conf, err := GetConfig()
	if err != nil {
		log.Printf("Error getting config [%s]", err)
		os.Exit(1)
	}

	var database db.DB
	switch conf.DBType {
	case "mem":
		database = db.NewMem()
	case "redis":
		redisOpts := &redis.Options{
			Addr:     conf.RedisHost,
			Password: conf.RedisPass,
			DB:       int(conf.RedisDB),
		}
		redisClient := redis.NewClient(redisOpts)
		database = db.NewRedis(redisClient)

	default:
		log.Printf("Error: no available DB type %s", conf.DBType)
		os.Exit(1)
	}

	router := mux.NewRouter()

	handler.NewCreateHandler(database).RegisterRoute(router)

	portStr := fmt.Sprintf(":%d", conf.Port)
	log.Printf("Serving on %s", portStr)
	log.Fatal(http.ListenAndServe(portStr, router))

}

func generateTLS() {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization:       []string{"Manning Publications Co."},
		OrganizationalUnit: []string{"Books"},
		CommonName:         "Go Web Programming",
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	certOut, _ := os.Create("cert.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	keyOut, _ := os.Create("key.pem")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
}

func serverTLS() {
	server := http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: nil,
	}
	server.ListenAndServeTLS("cert.pem", "key.pem")
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
