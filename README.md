# http_web-server
```
Web Api is develop using mux in golang
  url: localhost:8008
  Api end-point
      1. /signup 
          body container:
                {
                  "email":"{email_id}",
                  "name":"{user's_name}",
                  "password":"{user_password}"
                }
          return  message to confirm signup
      
      2. /login
          body Container:
                 body contains userid and password
                 {
                  "userid":"{user's_id}",
                  "password":"{password}"
                  }
                  
           return  A signed JSON Web token base64encodded
           
      
       3. /valid
            It's validate user's access token 
            token is passes in params
            i.e. /valid?token="{jwt}"
      ``` 
