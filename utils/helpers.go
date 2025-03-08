package utils

import(
    "os"
    "errors"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "log"
)

//GetCurrrentGoldPrice fetches the current market price of gold
//const mockGoldPrice = 1000.0

// The GetCurrentGoldPrice function returns the current gold price as a float64 value or an error.
func GetCurrentGoldPrice() (float64, error) {

    //Mocked function for testing 
    if  os.Getenv("TEST_ENV") == "true"{
        log.Println("Returning mocked gold price for testing...")
        return 2000.00, nil  // Mocked gold price for tests  
    }

    // Actual request to the API
    goldPriceAPI := "https://api.dukiapreciousmetals.co/api/price/products7"

    response, err := http.Get(goldPriceAPI)
    if err != nil {
        log.Printf("Error fetching gold price: %v", err)
        return 0, errors.New("failed to fetch gold price")
        }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    log.Println("Raw API Response:", string(body))
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return 0, errors.New("failed to read response body")
    }
    log.Println("Raw API Response:", string(body)) // Log the raw JSON response

    //Unmarshal into a slice (array of maps)
    var result []map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil{
        log.Printf("Error unmarshalling JSON: %v", err)
        return 0, errors.New("failed to unmarshall JSON")
    }
    
    // Ensure there is at least one item in the response body
    if len(result) == 0{
        return 0, errors.New("empty response from API")
    }

    // Extract the gold price from the response
    rate, ok := result[0]["ask_price"].(float64)
    if !ok {
        return 0, errors.New("invalid gold price format")
    }
    return rate, nil  // in USD
}

//SendNotificaton sends a notificaton to the user

func SendNotification(userID int, message string) error {
    log.Printf("Notification sent to User ID %d: %s\n", userID, message)
    return nil
}