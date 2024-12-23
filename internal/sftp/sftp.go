package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
)

const (
	homeDir = "home"
)

type SSHServer struct {
	Addr    string
	Port    int
	KeyPath string
}

func NewSSHServer(addr string, port int, keyPath string) *SSHServer {
	return &SSHServer{
		Addr:    addr,
		Port:    port,
		KeyPath: keyPath,
	}
}

// TODO: Upload a file to the sftp server, the code workds but with a local ssh server
func (s SSHServer) UploadFile() {
	// buho ssh server
	key, err := os.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}
	auth := ssh.PublicKeys(signer)

	// buho-sftp public key
	registeredPubKey, err := LoadRegisteredPublicKey("internal/sftp/pub_key")
	if err != nil {
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
		log.Fatal("Failed to dial: ", err)
	}

	defer conn.Close()

	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// user variable is a placeholder for the actual user requesting to upload a file
	var user = "USER"
	var sftpBasePath = fmt.Sprintf("%s/%s", homeDir, user)
	err = client.MkdirAll(sftpBasePath)
	if err != nil {
		log.Fatal(err)
	}

	// fileName is a placeholder for the actual file name
	fileName := "maple.JPG"
	// leave your mark
	f, err := client.Create(fmt.Sprintf("%s/%s", sftpBasePath, fileName))
	if err != nil {
		log.Fatal(err)
	}

	// read a file locally, however, this will be done with an incoming request,
	// most likely already []byte. This code wil be replaced with the actual file
	photo, err := os.Open("internal/sftp/maple-halloween.JPG")
	if err != nil {
		log.Fatal(err)
	}

	// write file to sftp server
	if _, err := f.ReadFrom(photo); err != nil {
		log.Fatal(err)
	}

	f.Close()

	// check it's there
	fi, err := client.Lstat(sftpBasePath + "/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	if fi == nil {
		log.Fatal("file not found")
	}
	fmt.Println("file uploaded successfully")
	client.Close()
	fmt.Println("client closed")
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
