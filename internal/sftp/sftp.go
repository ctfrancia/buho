package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"mime/multipart"
	"net"
	"os"
)

const (
	homeDir = "home"
)

type SSHServer struct {
	Addr              string
	Port              int
	PrivateKeyName    string
	AcceptedPublicKey string
}

func NewSSHServer(addr string, port int, pkn, apk string) *SSHServer {
	return &SSHServer{
		Addr:              addr,
		Port:              port,
		PrivateKeyName:    pkn,
		AcceptedPublicKey: apk,
	}
}

// UploadFile uploads a file to the sftp server
func (s SSHServer) UploadFile(file multipart.File, fileName, user string) (string, error) {
	// TODO: Clean up this function
	uploadPath := fmt.Sprintf("%s/%s", homeDir, user)
	remoteLocation := fmt.Sprintf("%s/%s", uploadPath, fileName)

	// buho ssh server
	key, err := os.ReadFile(s.PrivateKeyName)
	if err != nil {
		// TODO: error handle properly
		log.Fatal("Failed to load private key: ", err)
		return "", fmt.Errorf("Failed to load private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		// TODO: error handle properly
		log.Fatal("Failed to parse private key: ", err)
	}
	auth := ssh.PublicKeys(signer)

	// buho-sftp public key
	registeredPubKey, err := LoadRegisteredPublicKey("internal/sftp/pub_key")
	if err != nil {
		// TODO: error handle properly
		log.Fatal("Failed to load registered public key: ", err)
	}
	// ssh client config
	config := &ssh.ClientConfig{
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: HostKeyCb(registeredPubKey),
	}

	// connect to ssh server
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port), config)
	if err != nil {
		return "", fmt.Errorf("Failed to dial: %v", err)
	}

	defer conn.Close()

	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(conn)
	if err != nil {
		return "", fmt.Errorf("Error creating new sftp client: %v", err)
	}
	defer client.Close()

	// create the directory if it doesn't exist remotely
	err = client.MkdirAll(uploadPath)
	if err != nil {
		return "", fmt.Errorf("Err making all directories: %v", err)
	}

	// create a file on the remote server
	f, err := client.Create(fmt.Sprintf("%s/%s", uploadPath, fileName))
	if err != nil {
		return "", fmt.Errorf("Error creating file on remote server: %v", err)
	}

	// write file to sftp server
	if _, err := f.ReadFrom(file); err != nil {
		return "", fmt.Errorf("Error writing file to remote sftp server: %v", err)
	}

	f.Close()

	// check it's there
	fi, err := client.Lstat(remoteLocation)
	if err != nil {
		return "", fmt.Errorf("Error performing Lstat on path: %v", err)
	}
	if fi == nil {
		return "", fmt.Errorf("File not found on remote server")
	}
	client.Close()

	return remoteLocation, nil
}

func HostKeyCb(registeredKey ssh.PublicKey) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if string(key.Marshal()) == string(registeredKey.Marshal()) {
			return nil
		}

		return fmt.Errorf("host key mismatch")
	}
}

func LoadRegisteredPublicKey(path string) (ssh.PublicKey, error) {
	pubKeyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return pubKey, nil
}
