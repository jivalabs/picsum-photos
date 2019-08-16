package api

// There are two ways to send the authentication token to an API.
// You can include it as a query parameter, access_token=$token,
// or as an HTTP header Authorization: Bearer $token.
// The header method is recommended.

// # As a query parameter:
//curl -v https://cdn.contentful.com/spaces/cfexampleapi/entries?access_token=b4c0n73n7fu1
//curl -v https://cdn.jivalabs.com/spaces/cfexampleapi/entries?access_token=b4c0n73n7fu1

// # As a header:
//curl -v https://cdn.contentful.com/spaces/cfexampleapi/entries -H 'Authorization: Bearer b4c0n73n7fu1'

// if you fail to include a valid access token, you will receive an error message
// # Request
//curl https://cdn.contentful.com/spaces/cfexampleapi/entries?access_token=wrong
//
//# Response
//{
//  "sys": {
//    "type": "Error",
//    "id": "AccessTokenInvalid"
//  },
//  "message": "The access token you sent could not be found or is invalid.",
//  "requestId": "bcc-1808911724"
//}

// If you include a valid access token, but one that is not able to access a resource, you will receive a 404 error:
// # Request
//curl https://cdn.contentful.com/spaces/some_other_space/entries?access_token=b4c0n73n7fu1
//
//# Response
//{
//  "sys": {
//    "type": "Error",
//    "id": "NotFound"
//  },
//  "message": "The resource could not be found.",
//  "details": {
//    "sys": {
//      "type": "Space"
//    }
//  },
//  "requestId": "9f3-2148374087"
//}

// contentful Images API
