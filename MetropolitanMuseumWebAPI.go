package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var baseURL = "https://collectionapi.metmuseum.org/public/collection/v1/objects"
var objbaseURL = "https://collectionapi.metmuseum.org/public/collection/v1/objects/"

// Objects struct :json　→　構造体　に変換するようの構造体
type Objects struct {
	Total     int   `json:"total"`
	ObjectIDs []int `json:"objectIDs"`
}

// ObjectInfo struct : json　→　構造体　に変換するようの構造体
type ObjectInfo struct {
	ObjectID          int      `json:"objectID"`
	IsHighlight       bool     `json:"isHighlight"`
	AccessionNumber   string   `json:"accessionNumber"`
	AccessionYear     string   `json:"accessionYear"`
	IsPublicDomain    bool     `json:"isPublicDomain"`
	PrimaryImage      string   `json:"primaryImage"`      //URL to the primary image of an object in JPEG format
	PrimaryImageSmall string   `json:"primaryImageSmall"` //URL to the lower-res primary image of an object in JPEG format
	AdditionalImages  []string `json:"additionalImages"`
	Constituents      []struct {
		Role                   string `json:"role"`
		Name                   string `json:"name"`
		ConstituentULANURL     string `json:"constituentULAN_URL"`
		ConstituentWikidataURL string `json:"constituentWikidata_URL"`
		Gender                 string `json:"gender"`
	} `json:"constituents"`
	Department            string `json:"department"`
	ObjectName            string `json:"objectName"`
	Title                 string `json:"title"`
	Culture               string `json:"culture"`
	Period                string `json:"period"`
	Dynasty               string `json:"dynasty"`
	Reign                 string `json:"reign"`
	Portfolio             string `json:"portfolio"`
	ArtistRole            string `json:"artistRole"`
	ArtistPrefix          string `json:"artistPrefix"`
	ArtistDisplayName     string `json:"artistDisplayName"`
	ArtistDisplayBio      string `json:"artistDisplayBio"`
	ArtistSuffix          string `json:"artistSuffix"`
	ArtistAlphaSort       string `json:"artistAlphaSort"`
	ArtistNationality     string `json:"artistNationality"`
	ArtistBeginDate       string `json:"artistBeginDate"`
	ArtistEndDate         string `json:"artistEndDate"`
	ArtistGender          string `json:"artistGender"`
	ArtistWikidataURL     string `json:"artistWikidata_URL"`
	ArtistULANURL         string `json:"artistULAN_URL"`
	ObjectDate            string `json:"objectDate"`
	ObjectBeginDate       int    `json:"objectBeginDate"`
	ObjectEndDate         int    `json:"objectEndDate"`
	Medium                string `json:"medium"`
	Dimensions            string `json:"dimensions"`
	CreditLine            string `json:"creditLine"`
	GeographyType         string `json:"geographyType"`
	City                  string `json:"city"`
	State                 string `json:"state"`
	County                string `json:"county"`
	Country               string `json:"country"`
	Region                string `json:"region"`
	Subregion             string `json:"subregion"`
	Locale                string `json:"locale"`
	Locus                 string `json:"locus"`
	Excavation            string `json:"excavation"`
	River                 string `json:"river"`
	Classification        string `json:"classification"`
	RightsAndReproduction string `json:"rightsAndReproduction"`
	LinkResource          string `json:"linkResource"`
	MetadataDate          string `json:"metadataDate"`
	Repository            string `json:"repository"`
	ObjectURL             string `json:"objectURL"`
	Tags                  []struct {
		Term   string `json:"term"`
		AATURL string `json:"AAT_URL"`
	} `json:"tags"`
	ObjectWikidataURL string `json:"objectWikidata_URL"`
	IsTimelineWork    bool   `json:"isTimelineWork"`
}

// ObjIDStartNum 画像取得開始するリスト
const ObjIDStartNum = 40000

func main() {
	/*** Object List 取得 ***/
	resp, err := http.Get(baseURL)
	if err != nil {
		panic(err)
	}
	obj := new(Objects)
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(byteArray, &obj); err != nil {
		panic(err)
	}

	fmt.Println("***** 1 *****")

	/* パブリックドメインデータがあるobjIDの画像データを取得 */
	pdObjID := 0 /* パブリックドメインObjID格納変数初期化 */
	objinfo := new(ObjectInfo)

	/*** 読み出しファイルオープン ***/
	// O_RDONLY:読み出しモード開く
	rfile, err := os.OpenFile("./in/objlist.csv", os.O_RDONLY, 0600) //ファイルに追記する場合
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer rfile.Close()

	reader := csv.NewReader(rfile)
	reader.LazyQuotes = true

	fmt.Println("***** 2 *****")
	fmt.Printf("ObjIDStartNum:%d \n", 0)
	fmt.Printf("ObjID:%d check\n", obj.ObjectIDs[0])

	/*** 書き込みファイルオープン ***/
	// O_WRONLY:書き込みモード開く, O_CREATE:無かったらファイルを作成
	wfile, err := os.OpenFile("./out/imageList.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600) //ファイルに追記する場合
	if err != nil {
		log.Fatal("Error:", err)
	}
	writer := csv.NewWriter(wfile)

	//項目名
	writer.Write([]string{"ObjectID", "ArtistName", "Classification", "ImageURL", "ThumnailURL"})

	// 読み出しファイルのObjectListの数だけループ
	row, err := reader.Read()
	fmt.Printf("row: %x\n", row)

	for ; err != io.EOF; row, err = reader.Read() {
		/* ファイル読み出しエラー時 */
		if err != io.EOF && err != nil {
			fmt.Fprintln(os.Stderr, "CSV形式のデータを正しく読み取れませんでした。")
		}

		// objectIDをファイルから取得
		objIDText, err := reader.Read()
		if err != nil {
			log.Fatal(err)
			break
		}
		objID := objIDText[0]

		objURL := objbaseURL + objID
		resp3, err := http.Get(objURL)
		if err != nil {
			panic(err)
		}
		defer resp3.Body.Close()
		byteArray, err = ioutil.ReadAll(resp3.Body)
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal(byteArray, &objinfo); err != nil {
			fmt.Printf("ObjID:%x Error\n", objID)
			panic(err)
		}

		/* 取得したobjIDのpublicDomain情報確認 */
		isPD := objinfo.IsPublicDomain
		classifiCation := objinfo.Classification
		if isPD == true && classifiCation == "Paintings" {
			pdObjID, _ = strconv.Atoi(objID)
			artistName, imgURL, smallImgURL := SaveArtList(objinfo, pdObjID)
			if artistName != "" {
				//画家の名前情報があるデータのみ出力
				writer.Write([]string{strconv.Itoa(pdObjID), artistName, classifiCation, imgURL, smallImgURL})

				// サムネイル画像保存
				SaveImage(pdObjID, smallImgURL)

				fmt.Printf("ObjID:%d Success\n", pdObjID)
				time.Sleep(time.Millisecond * 100) // 100ミリ秒 Sleep
			}
		} else {

		}
	}

	/* ファイル書き出し、クローズ処理 */
	writer.Flush()      //ファイルに書き出し
	defer wfile.Close() //ファイルクローズ
}

// SaveImage function
func SaveImage(objectID int, imgURL string) {
	/* PublicDomain画像取得 */
	resImg, err := http.Get(imgURL)
	if err != nil {
		panic(err)
	}
	defer resImg.Body.Close()

	//jpgファイル作成＆保存
	filename := "save" + strconv.Itoa(objectID) + ".jpg"
	file, err := os.Create("./images/" + filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(file, resImg.Body)

}

// SaveArtList function
func SaveArtList(objInfo *ObjectInfo, pdObjectID int) (string, string, string) {
	artistName := objInfo.ArtistDisplayName
	imgURL := objInfo.PrimaryImage
	smallImgURL := objInfo.PrimaryImageSmall

	return artistName, imgURL, smallImgURL
}
