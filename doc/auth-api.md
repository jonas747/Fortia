#Auth api methods

###No atuthentication required
    
####/login
Method: POST
Returns a session cookie that expires after 1h inactivity
Or an error
Body:

    {
        "user": "username to login",
        "pw": "password to login with"
    }

####/register
Method: POST
Returns an error if something went wrong. {"ok": true} if everything went okay
Body:

    {
        "username": "username to register",
        "pw": "password to register",
        "email": "email to register"
    }

####/worlds
Method: GET
Returns info about all worlds or the specified world
optional params: world

###Authentication required

####/me
Method: GET
Returns info about yourself