package export

import (
    "encoding/json"
 
    "github.com/caltechlibrary/crossrefapi"
)

func Json(r *crossrefapi.Message) (string, error) {
    jsonBytes, err := json.MarshalIndent(*r, "", "  ")
    if err != nil {
        return "", err
    }
    return string(jsonBytes), nil
}
