package chaincode

const (
	S1      = " \n{\n\"type\": \"GOLANG\",\n\t\"chaincodeID\": {\n\t\t"
	S2      = "\"\n},\n\"ctorMsg\": {\n\t\"function\": \""
	S3      = "\",\n\t\"args\": ["
	S4NOSEC = "]\n} }"
	S4      = "]\n},\n\"secureContext\":\""
	S5      = "\"\n}"
	IQSTART = "\n{\"chaincodeSpec\":"
	IQEND   = "\n}"

	escQuote          = "\""
	openBrace         = "{"
	closeBrace        = "}"
	comma             = ","
	colon             = ":"
	newLine           = "\n"
	RegisterJsonPart1 = newLine + openBrace + newLine + escQuote + "enrollId" + escQuote + colon + escQuote
	RegisterJsonPart2 = escQuote + comma + newLine + escQuote + "enrollSecret" + escQuote + colon + escQuote
	RegisterJsonPart3 = escQuote + newLine + closeBrace + newLine
)

// Use this counter to generate a distinct value for "id" in the REST API payload
// Initialize in LoadNetwork()
// Increment in genPayLoadForChaincode()
var (
	PostChaincodeCount int64 = 1
)
