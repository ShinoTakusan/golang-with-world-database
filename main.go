package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

type CityList struct {
	Name string `json:"name,omitempty"  db:"Name"`
}

type Country struct {
	Code string `json:"code,omitempty"  db:"Code"`
	Name string `json:"name,omitempty"  db:"Name"`
}

var (
	db *sqlx.DB
)

func main() {
	// データベースにアクセス
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	//echoの設定
	e := echo.New()

	//セッションを覚えておく場所？？
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.POST("/login", postLoginHandler)
	e.POST("/signup", postSignUpHandler)

	withLogin := e.Group("") //何もないグループの宣言(linuxの権限グループ的な)
	withLogin.Use(checkLogin)
	withLogin.GET("/cities/:cityName", getCityInfoHandler)
	withLogin.GET("/whoami", getWhoAmIHandler)
	withLogin.GET("/country", getCountryInfoHandler)
	withLogin.GET("/citylist/:countryName",getCityListHandler)

	e.Start(":12200")
}

type LoginRequestBody struct {
	Username string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

type User struct {
	Username   string `json:"username,omitempty"  db:"Username"`
	HashedPass string `json:"-"  db:"HashedPass"`
}

type Me struct {
	Username string `json:"username,omitempty" db:"username"`
}

//User登録を行う関数
func postSignUpHandler(c echo.Context) error {
	//reqに認証情報をいれる。
	req := LoginRequestBody{}
	c.Bind(&req)

	//なんも入ってないとき
	// もう少し真面目にバリデーションするべき
	if req.Password == "" || req.Username == "" {
		// エラーは真面目に返すべき
		return c.String(http.StatusBadRequest, "項目が空です")
	}

	//本当に入ってないのかを検査
	//はいっていたら、ハッシュ化して保存
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	// ユーザーの存在チェック
	var count int
	// ユーザーの人数をとる
	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE Username=?", req.Username)
	// データーベースエラー
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	//すでにUSERがいる場合は、存在しているので駄目ですとする。
	if count > 0 {
		return c.String(http.StatusConflict, "ユーザーが既に存在しています")
	}

	//データベースに情報を追加、駄目な場合はinternal error
	_, err = db.Exec("INSERT INTO users (Username, HashedPass) VALUES (?, ?)", req.Username, hashedPass)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}

// ログインに関する関数
func postLoginHandler(c echo.Context) error {
	//ログイン情報をバインド
	req := LoginRequestBody{}
	c.Bind(&req)

	//ユーザーネームを参照し、パスワードとユーザーネームをとってくる
	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE username=?", req.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	//パスワードをハッシュ化して比較
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
	if err != nil {
		//ハッシュが一致しない
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.NoContent(http.StatusForbidden)

		} else {
			//ハッシュ化ができない
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	//セッションを保存する
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	sess.Values["userName"] = req.Username
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}

//middleware (ハンドラーを返す関数)
func checkLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//セッションの取得
		sess, err := session.Get("sessions", c)
		//セッションが取れなかった…
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}
		//セッションがないときは、ログインしてない。
		if sess.Values["userName"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		//ログインしていることを確認したので、cにいれる。
		c.Set("userName", sess.Values["userName"].(string))
		fmt.Println("login !!", c.Get("userName"))
		return next(c)
	}
}

//Paramで与えられたものを返す関数
func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")

	city := City{}
	db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName)
	if city.Name == "" {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, city)
}

func getWhoAmIHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Me{
		Username: c.Get("userName").(string),
	})
}

func getCountryInfoHandler(c echo.Context) error {
	country := []Country{}
	db.Select(&country, "SELECT Code, Name FROM country")
	return c.JSON(http.StatusOK, country)
}

func getCityListHandler(c echo.Context) error {
	countryName := c.Param("countryName")

	citylist := []CityList{}
	db.Select(&citylist,
		"select city.Name from city join country on CountryCode = Code where country.Name = ?",
		countryName)
	if len(citylist) == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, citylist)
}
