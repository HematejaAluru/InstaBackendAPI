# InstaBackendAPI
# Quick Tour:
To use InstaBackendAPI, please follow the steps:
  1. Unzip the given folder
  2. Make sure you have golang setup in local
  3. Make sure you have MongoDB setup
  4. To run the API please use these commands "go build main.go" and "./main" inside the FinalInstaApi
  5. You can use postman to call the APIs


# Functionalities:
  1. FinalInstaApi has 5 different Apis - "createUser","getUser","createPost","getPost","getAllPosts".
  2. getAllPosts also has pagination. It has limit and offset parameters, Frontend can provide the required limit and offset so that we can provide seamless experience to user.
  3. For storing user's password safely in the MongoDB, SHA256 algorithm is used for hashing.
  4. For tesing the FinalInstaApi, there is a folder named "TestingInstaapi" which contains unit tests for all Apis
