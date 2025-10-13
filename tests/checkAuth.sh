curl -X GET http://localhost:8081/checkAuth
curl -X GET http://localhost:8081/checkAuth \
  -H "userID: 123" \
  -H "userKEY: valid_key_123"
