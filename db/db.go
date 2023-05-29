package db

import (
	"github.com/boltdb/bolt"
	"log"
	"strings"
	"time"
)

var r *bolt.DB

const USER_BUCKET_NAME = "$:USER"
const DOMAIN_OWNER = "$:DOMAIN:"

var domain = []string{"dev", "test", "prod", "db", "io", "api", "app"}

type User struct {
	Name           string `json:"username"`
	Password       string `json:"password"`
	RetypePassword string `json:"retype_password"`
}

func Init() *bolt.DB {
	if r == nil {
		database, err := bolt.Open("./rapid.dat", 0600, nil)
		if err != nil {
			log.Fatal("Database couldn't open!")
		}
		r = database
	}
	return r
}

func LoadDefault() {
	defaultHost := "rapid"
	user := User{Name: defaultHost, Password: defaultHost}
	AddNewUser(user)
	for _, d := range domain {
		r.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte(d))
			if err != nil {
				log.Fatal("Create bucket failed for ", d)
			}
			bucket.Put([]byte(defaultHost), []byte("69.69.69.69"))
			return nil
		})
		MarkDomainForUser(defaultHost+"."+d, user.Name)
	}

}

func ListUsableDomains(name string) []string {
	var list []string
	log.Println(name)
	r.View(func(tx *bolt.Tx) error {
		for _, d := range domain {
			domainBucket := tx.Bucket([]byte(d))
			log.Println(d + "  " + name)
			exist := domainBucket.Get([]byte(name))
			log.Println(exist)
			if exist == nil {
				list = append(list, name+"."+d)
			}
		}
		return nil
	})
	return list
}

func ExistSameUsername(user User) bool {
	rst := false
	r.View(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(USER_BUCKET_NAME))
		if userBucket != nil {
			rst = nil != userBucket.Get([]byte(user.Name))
		}
		return nil
	})
	return rst
}

func ValidateUser(user User) bool {
	rst := false
	r.View(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(USER_BUCKET_NAME))
		if userBucket != nil {
			password := userBucket.Get([]byte(user.Name))
			rst = user.Password == string(password)
		}
		return nil
	})
	return rst
}

func AddNewUser(user User) bool {
	rst := false
	r.Update(func(tx *bolt.Tx) error {
		userBucket, err := tx.CreateBucketIfNotExists([]byte(USER_BUCKET_NAME))
		if err != nil {
			log.Fatal("Couldn't create USER_BUCKET for user")
		}
		err = userBucket.Put([]byte(user.Name), []byte(user.Password))
		_, err = tx.CreateBucketIfNotExists([]byte(DOMAIN_OWNER + user.Name))
		if err != nil {
			log.Fatal("Couldn't create DOMAIN_OWNER for user")
		}
		rst = err == nil
		return nil
	})
	return rst
}

func Query(name string) []byte {
	split := strings.Split(name, ".")
	var resp []byte
	r.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(split[len(split)-2]))
		if bucket == nil {
			return nil
		}
		resp = bucket.Get([]byte(split[len(split)-3]))
		return nil
	})
	return resp
}

func ExistDomain(name string) bool {
	split := strings.Split(name, ".")
	rst := true
	r.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(split[len(split)-1]))
		if bucket == nil {
			rst = false
			return nil
		}
		rst = bucket.Get([]byte(split[len(split)-2])) != nil
		return nil
	})
	return rst
}

func DeleteDomainWithIpv4(name string) bool {
	split := strings.Split(name, ".")
	bucketName := split[len(split)-1]
	domainName := split[len(split)-2]
	rst := true
	r.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		bucket.Delete([]byte(domainName))
		return nil
	})
	return rst
}

func AddDomainWithIpv4(name string, ipv4 string) bool {
	split := strings.Split(name, ".")
	domainName := split[len(split)-2]
	bucketName := split[len(split)-1]
	rst := true
	r.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		bucket.Put([]byte(domainName), []byte(ipv4))
		return nil
	})
	return rst
}

func OwnDomain(domainName string, username string) bool {
	var own = false
	r.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DOMAIN_OWNER + username))
		own = bucket.Get([]byte(domainName)) != nil
		return nil
	})
	return own
}

func DomainListOfUser(username string) []string {
	var domains []string
	r.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DOMAIN_OWNER + username))

		return bucket.ForEach(func(k, v []byte) error {
			domains = append(domains, string(k))
			return nil
		})
	})
	return domains
}

func MarkDomainForUser(domainName string, username string) {
	r.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DOMAIN_OWNER + username))
		bucket.Put([]byte(domainName), []byte(string(time.Now().Unix())))
		return nil
	})
}
func UnMarkDomainForUser(domainName string, username string) {
	r.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DOMAIN_OWNER + username))
		bucket.Delete([]byte(domainName))
		return nil
	})
}
