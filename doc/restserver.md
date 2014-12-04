## Rest server documentation

A Handler

    type RestHandler struct {
        Handler         RestHandlerFunc
        Method          string                  // The metho ex: GET, PUT, PATCH etc..
        RequiredParams  []string                // Required url parameters
        OptionalParams  []string                // Optional Url parameters
        Path            string                  // The path this handler takes action upon
        BodyType        reflect.Type            // The type of the body
        BodyRequired    bool                    // Wther a body is required or not
        AdditionalData  map[string]interface{} // Additional data
        Middleware []RestHandlerFunc
    }

    