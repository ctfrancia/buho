package sftp

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"log"
	"mime/multipart"
	"net"

	bAuth "github.com/ctfrancia/buho/internal/auth"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	homeDir = "home"
)

type SSHServer struct {
	Addr           string
	Port           int
	PrivateKeyPath string
	PublicKeyPath  string
	SFTPKeyPath    string
}

func NewSSHServer(addr string, port int, pubKeyPath, privKeyPath string) *SSHServer {
	return &SSHServer{
		Addr:           addr,
		Port:           port,
		PublicKeyPath:  pubKeyPath,
		PrivateKeyPath: privKeyPath,
		SFTPKeyPath:    "internal/keys/sftp/buho-sftp.pem",
	}
}

// UploadFile uploads a file to the sftp server
func (s SSHServer) UploadFile(file multipart.File, fileName, website string) (string, error) {
	uploadPath := fmt.Sprintf("%s/%s", homeDir, website)
	remoteLocation := fmt.Sprintf("%s/%s", uploadPath, fileName)

	// buho ssh server
	buhoPrivKey, err := bAuth.ParseED25519PrivateKey(s.PrivateKeyPath)
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
		return "", fmt.Errorf("Failed to load private key: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(buhoPrivKey)
	if err != nil {
		return "", fmt.Errorf("failed to create signer: %v", err)
	}

	auth := ssh.PublicKeys(signer)

	fmt.Println("public key path: ", s.SFTPKeyPath)
	registeredPubKey, err := bAuth.ParseED25519PublicKey(s.SFTPKeyPath)
	if err != nil {
		fmt.Println("error loading registered public key ----------- ", s.PublicKeyPath)
		log.Fatal(err)
	}
	// Add the private key to the SSH client config's authentication methods
	config := &ssh.ClientConfig{
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: HostKeyCb(registeredPubKey),
	}

	// connect to ssh server
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port), config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %v", err)
	}

	defer conn.Close()

	// TODO: START HERE
	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(conn)
	if err != nil {
		return "", fmt.Errorf("error creating new sftp client: %v", err)
	}
	defer client.Close()

	// create the directory if it doesn't exist remotely
	err = client.MkdirAll(uploadPath)
	if err != nil {
		return "", fmt.Errorf("err making all directories: %v", err)
	}

	// create a file on the remote server
	f, err := client.Create(fmt.Sprintf("%s/%s", uploadPath, fileName))
	if err != nil {
		return "", fmt.Errorf("error creating file on remote server: %v", err)
	}

	// write file to sftp server
	if _, err := f.ReadFrom(file); err != nil {
		return "", fmt.Errorf("error writing file to remote sftp server: %v", err)
	}

	f.Close()

	// check it's there
	fi, err := client.Lstat(remoteLocation)
	if err != nil {
		return "", fmt.Errorf("error performing Lstat on path: %v", err)
	}
	if fi == nil {
		return "", fmt.Errorf("file not found on remote server")
	}
	client.Close()

	return remoteLocation, nil
}

func HostKeyCb(registeredKey ed25519.PublicKey) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		// Log incoming information for debugging
		log.Printf("Host Key Callback Debug:")
		log.Printf("Hostname: %s", hostname)
		log.Printf("Remote Address: %s", remote)
		log.Printf("Incoming Key Type: %T", key)

		// Convert registered ED25519 public key to SSH public key
		sshPubKey, err := ssh.NewPublicKey(registeredKey)
		if err != nil {
			return fmt.Errorf("failed to convert registered public key: %w", err)
		}

		// Log key details for comparison
		log.Printf("Registered Key (SSH): %x", sshPubKey.Marshal())
		log.Printf("Incoming Key: %x", key.Marshal())

		// Perform detailed key comparison
		registeredKeyBytes := sshPubKey.Marshal()
		incomingKeyBytes := key.Marshal()

		if bytes.Equal(registeredKeyBytes, incomingKeyBytes) {
			return nil
		}

		// If keys don't match, provide detailed mismatch information
		log.Printf("Host Key Mismatch:")
		log.Printf("Registered Key Length: %d", len(registeredKeyBytes))
		log.Printf("Incoming Key Length: %d", len(incomingKeyBytes))

		return fmt.Errorf("host key mismatch: registered key does not match incoming key")
	}
}
