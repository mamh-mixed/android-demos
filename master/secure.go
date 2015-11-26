package master

import "github.com/CardInfoLink/quickpay/security"

var privateKeyForBrowser = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDF4xejo7F1JVPU555mG6Kei8XU2bT+V0Y+DaxzoBChaxYGOtlk
f6vCh3y6Op/3OWdZAG8W17S3w9V7Skw0PvFvqqqc8JLlnr9/zDKoit5X17VHX8Ky
3jdl7Ll2h3MFghAbzcf0P7CRGxgpTm+lqsPQXETzDEBEqXeE7Q7WeseaHQIDAQAB
AoGBALRCaH89FuLibNoNX0IePGV2Z3C8HF5vu+G87PGqxlt0Q+zK4MrmbdzXNKwj
ySIYXWc6uPcy6UFYl/gmNwKEr8JvNzt7X/ZxPz1KUC7lge3yHfqdVBvzKb2gJxiq
330qxIRDnNbbu8wFKroY/lToKYaPDbyJmgzwNpgZNQ2NIZLhAkEA+5crL4iyskiN
+ZDWTIikM0JGNz5mkR+FZzO2n7c51nQnvHivu/Pyt4OVLtTJNYWJntj5s4j+16eS
nOhZzo4+iQJBAMla99sN55nWSjvtI1AWrSmpbmX3s0juRz/oy2KB/ImO4Oubhs6w
/GYdo/phtokWTIAxopUWL1pJHKsUnWqfefUCQQDeNxYIxQd4msbzoC73qFTHhYj5
MF9tXNb6YV2zUiV+uleCi2JEc2J1Hn58v6r8X/c+20wpfB4DIlpHxp3T6CVpAkEA
guJGynU3XqAUkO+MTLrwxGwF/vIL8BQy7C/+RIIKDcB6I6xs7F3PMvGBbXeml2WP
RKT+8boB/cYYhHxZ9rzDIQJAWYu4Erm8Oaed1NnXlBkLiz6jq8DVc4x2zpAK+VbZ
0E0blJfkbmpNdOzyYLr+LfoQ7/99EiSOfg77ZkLeHby9QA==
-----END RSA PRIVATE KEY-----
	`)

func rsaDecryptFromBrowser(cipertext string) (plaintext string, err error) {
	plainData, err := security.RSADecryptBase64(cipertext, privateKeyForBrowser)
	if err != nil {
		return "", err
	}

	// n := bytes.IndexByte(plainData, 0)
	plaintext = string(plainData)
	return
}
