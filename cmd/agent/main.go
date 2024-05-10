package main

import (
	"github.com/bbquite/mca-server/internal/handlers"
	"log"
)

func main() {
	if err := handlers.AgentRun(); err != nil {
		log.Fatal(err)
	}
}

//st := new(runtime.MemStats)
//repo := &storage.Repository{}
//log.Print(repo)
//repo.AddGaugeItem(storage.MemStorageGauge{Key: "test", Value: 1.25})
//
//log.Print(repo)

//memHandler := &storage.MemStorageHandler{
//	storage: &storage.MemStorage{
//		items: make(map[string]storage.MemStorageItem),
//	},
//}
//
//memHandler.CreateMemItem("Alloc", string(st.Alloc))

//for {
//	st := new(runtime.MemStats)
//	runtime.ReadMemStats(st)
//	test := st.Alloc
//	log.Println(test)
//	time.Sleep(2 * time.Second)
//}

// response, err := http.Get("https://practicum.yandex.ru")
// if err != nil {
// 	fmt.Println(err)
// }
// fmt.Printf("Status Code: %d\r\n", response.StatusCode)
// for k, v := range response.Header {
// 	// заголовок может иметь несколько значений,
// 	// но для простоты запросим только первое
// 	fmt.Printf("%s: %v\r\n", k, v[0])
// }
