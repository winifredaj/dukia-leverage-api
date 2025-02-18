package utils

import(
    "errors"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "log"
)

//GetCurrrentGoldPrice fetches the current market price of gold
const goldPriceAPI = "https://api.metalsapi.com/latest?access_key=YOUR_API_KEY&base=USD&symbols=XAU"

func GetCurrentGoldPrice() (float64, error) {
    response, err := http.Get(goldPriceAPI)
    if err != nil {
        log.Printf("Error fetching gold price: %v", err)
        return 0, errors.New("failed to fechgold price")
        }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return 0, errors.New("failed to read response body")
    }

    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil{
        log.Printf("Error unmarshalling JSON: %v", err)
        return 0, errors.New("failed to unmarshall JSON")
    }

    rate, ok := result["rates"].(map[string]interface{})["XAU"].(float64)
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