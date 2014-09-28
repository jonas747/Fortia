#Auth api methods

####/Login
Method: POST
Returns a session cookie that expires after 1h inactivity
Or an error
Body:

    {
        "user": "username to login",
        "pw": "password to login with"
    }

####/Register
Method: POST
Returns an error if something went wrong. {"ok": true} if everything went okay
Body:

    {
        "username": "username to register",
        "pw": "password to register",
        "email": "email to register"
    }

####/GetInfo
Method: GET
Returns info about yourself
params: user