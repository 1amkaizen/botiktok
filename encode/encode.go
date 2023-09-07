package encode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Fungsi untuk menghapus markup HTML dari teks
func stripHTML(html string) string {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	return doc.Text()
}

// Fungsi untuk membaca isi file JSON
func ReadJSONFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Fungsi untuk mencari produk berdasarkan kata kunci
func SearchProduct(filePath, keyword string) ([]string, error) {
	// Membaca file JSON
	data, err := ReadJSONFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Gagal membaca file: %v", err)
	}

	// Mendekode JSON
	jsonData, err := DecodeJSON(data)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengurai JSON: %v", err)
	}

	// Mendapatkan daftar produk dari JSON
	products := jsonData["data"].(map[string]interface{})["products"].([]interface{})

	// Membuat slice untuk produk yang sesuai
	matchingProducts := make([]string, 0)

	// Mencari produk yang sesuai dengan kata kunci
	for _, product := range products {
		productName := product.(map[string]interface{})["product_name"].(string)
		if strings.Contains(productName, keyword) {
			matchingProducts = append(matchingProducts, productName)
		}
	}

	return matchingProducts, nil
}

// Fungsi untuk mendekode JSON
func DecodeJSON(data []byte) (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}

// Fungsi GetDesc
func GetDesc(filePath string) (string, error) {
	// Membaca file JSON
	data, err := ReadJSONFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Gagal membaca file: %v", err)
	}

	// Mendekode JSON
	jsonData, err := DecodeJSON(data)
	if err != nil {
		return "", fmt.Errorf("Gagal mengurai JSON: %v", err)
	}

	// Mendapatkan deskripsi HTML dari JSON
	descriptionHTML := jsonData["data"].(map[string]interface{})["description"].(string)
	// Menghapus markup HTML
	descriptionText := stripHTML(descriptionHTML)

	return descriptionText, nil
}

// Fungsi untuk mendapatkan kategori dari JSON
func GetCategories(filePath string) ([]string, error) {
	// Membaca file JSON
	data, err := ReadJSONFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Gagal membaca file: %v", err)
	}

	// Mendekode JSON
	jsonData, err := DecodeJSON(data)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengurai JSON: %v", err)
	}

	// Mendapatkan kategori dari JSON
	categories := jsonData["data"].(map[string]interface{})["category_list"].([]interface{})
	categoryList := make([]string, 0)
	for _, category := range categories {
		categoryData := category.(map[string]interface{})
		categoryStr, ok := categoryData["local_display_name"].(string)
		if ok {
			categoryList = append(categoryList, categoryStr)
		}
	}

	return categoryList, nil
}

// Fungsi untuk mendapatkan informasi gambar dari JSON
func GetImages(filePath string) ([]string, error) {
	// Membaca file JSON
	data, err := ReadJSONFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Gagal membaca file: %v", err)
	}

	// Mendekode JSON
	jsonData, err := DecodeJSON(data)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengurai JSON: %v", err)
	}

	// Mendapatkan informasi gambar dari JSON
	images := jsonData["data"].(map[string]interface{})["images"].([]interface{})
	imageInfoList := make([]string, 0)
	for _, image := range images {
		imageData := image.(map[string]interface{})
		imageStr := fmt.Sprintf("ID: %s, Lebar: %d, Tinggi: %d", imageData["id"], int(imageData["width"].(float64)), int(imageData["height"].(float64)))
		imageInfoList = append(imageInfoList, imageStr)
	}

	return imageInfoList, nil
}

// Fungsi untuk mendapatkan URL dari JSON
// Fungsi untuk mendapatkan URL yang unik dari JSON
func GetUniqueURLs(filePath string) ([]string, error) {
	// Membaca file JSON
	data, err := ReadJSONFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Gagal membaca file: %v", err)
	}

	// Mendekode JSON
	jsonData, err := DecodeJSON(data)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengurai JSON: %v", err)
	}

	// Mendapatkan URL dari JSON
	urls := jsonData["data"].(map[string]interface{})["images"].([]interface{})
	urlSet := make(map[string]struct{}) // Menggunakan map untuk menghindari duplikat
	uniqueURLs := make([]string, 0)

	for _, urlData := range urls {
		url := urlData.(map[string]interface{})["thumb_url_list"].([]interface{})
		for _, urlValue := range url {
			urlStr, ok := urlValue.(string)
			if ok && !isURLInSet(urlStr, urlSet) {
				urlSet[urlStr] = struct{}{}
				uniqueURLs = append(uniqueURLs, urlStr)
			}
		}
	}

	return uniqueURLs, nil
}

// Fungsi untuk memeriksa apakah URL sudah ada dalam set
func isURLInSet(url string, urlSet map[string]struct{}) bool {
	_, exists := urlSet[url]
	return exists
}

func GetProductName(filePath string) (string, error) {
	// Membaca file JSON
	data, err := ReadJSONFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Gagal membaca file: %v", err)
	}

	// Mendekode JSON
	jsonData, err := DecodeJSON(data)
	if err != nil {
		return "", fmt.Errorf("Gagal mengurai JSON: %v", err)
	}

	// Mendapatkan deskripsi HTML dari JSON
	nameHTML := jsonData["data"].(map[string]interface{})["product_name"].(string)
	// Menghapus markup HTML
	nameText := stripHTML(nameHTML)

	return nameText, nil
}

// Fungsi untuk mengambil waktu pembuatan dari JSON
func GetCreateTime(filePath string) (int64, error) {
	// Membaca file JSON
	data, err := ReadJSONFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("Gagal membaca file: %v", err)
	}

	// Mendekode JSON
	jsonData, err := DecodeJSON(data)
	if err != nil {
		return 0, fmt.Errorf("Gagal mengurai JSON: %v", err)
	}

	// Mendapatkan waktu pembuatan dari JSON
	createTime := jsonData["data"].(map[string]interface{})["create_time"].(int64)

	return createTime, nil
}
