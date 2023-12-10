package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type InputData struct {
	ToSort [][]int `json:"to_sort"`
}

type ProcessedData struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNs       int64   `json:"time_ns"`
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/",fileServer)

	http.HandleFunc("/process-single", ProcessSingleHandler)
	http.HandleFunc("/process-concurrent", ProcessConcurrentHandler)

	fmt.Println("Server listening on :8000...")
	http.ListenAndServe(":8000", nil)
}

func ProcessSingleHandler(w http.ResponseWriter, r *http.Request) {
	var inputData InputData
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := processSingle(inputData.ToSort)
	duration := time.Since(startTime)

	response := ProcessedData{
		SortedArrays: sortedArrays,
		TimeNs:       duration.Nanoseconds(),
	}

	sendJSONResponse(w, response)
}

func ProcessConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	var inputData InputData
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := processConcurrent(inputData.ToSort)
	duration := time.Since(startTime)

	response := ProcessedData{
		SortedArrays: sortedArrays,
		TimeNs:       duration.Nanoseconds(),
	}

	sendJSONResponse(w, response)
}

// i am shivanand kerur as i used quick sort for sequential searches because it is faster and most optimized approach as far as i known
func processSingle(toSort [][]int) [][]int {
	sortedArrays := make([][]int, len(toSort))
	for i, arr := range toSort {
		sortedSubArr := make([]int, len(arr))
		copy(sortedSubArr, arr)
		quicksort(sortedSubArr)
		sortedArrays[i] = sortedSubArr
	}
	return sortedArrays
}

// here i used parallel quick sort for sequential searches because it is faster and most optimized approach i was thinking about using merge sort but is less 
// efficient than quick sort and many deadlocks and panic mode quite triggered more when i want use merge sort
func processConcurrent(toSort [][]int) [][]int {
	var wg sync.WaitGroup
	wg.Add(len(toSort))
	out := make(chan []int, len(toSort))

	for _, arr := range toSort {
		go func(arr []int) {
			defer wg.Done()
			sortedSubArr := make([]int, len(arr))
			copy(sortedSubArr, arr)
			quicksort(sortedSubArr)
			out <- sortedSubArr
		}(arr)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	sortedArrays := make([][]int, 0, len(toSort))
	for sortedSubArr := range out {
		sortedArrays = append(sortedArrays, sortedSubArr)
	}

	sort.Slice(sortedArrays, func(i, j int) bool {
		return sortedArrays[i][0] < sortedArrays[j][0]
	})

	return sortedArrays
}


func quicksort(arr []int) {
	if len(arr) <= 1 {
		return
	}

	pivot := arr[len(arr)/2]
	left, right := partition(arr, pivot)

	quicksort(left)
	quicksort(right)

	copy(arr, append(append(left, pivot), right...))
}

func partition(arr []int, pivot int) (left, right []int) {
	for _, v := range arr {
		if v < pivot {
			left = append(left, v)
		} else if v > pivot {
			right = append(right, v)
		}
	}
	return
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
