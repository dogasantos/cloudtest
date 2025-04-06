package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Testa listagem de blobs em um container via acesso anônimo.
func testBlobList(account, container string) error {
	// URL para listar os blobs do container.
	// A operação "List Blobs" é feita com: ?restype=container&comp=list
	url := fmt.Sprintf("https://%s.blob.core.windows.net/%s?restype=container&comp=list", account, container)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("Listagem anônima do container %q retornou: %s\n", container, resp.Status)
	return nil
}

// Testa upload de blob via acesso anônimo.
func testBlobUpload(account, container, blobName string) error {
	// URL para criar um blob
	url := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", account, container, blobName)
	req, err := http.NewRequest("PUT", url, strings.NewReader("conteúdo de teste"))
	if err != nil {
		return err
	}
	// Cabeçalho obrigatório para operação de Put Blob.
	req.Header.Set("x-ms-blob-type", "BlockBlob")
	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("Upload anônimo do blob %q no container %q retornou: %s\n", blobName, container, resp.Status)
	return nil
}

// Testa listagem de arquivos em um file share via acesso anônimo.
func testFileList(account, share string) error {
	// URL para listar o diretório raiz do file share.
	url := fmt.Sprintf("https://%s.file.core.windows.net/%s?restype=directory&comp=list", account, share)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("Listagem anônima do file share %q retornou: %s\n", share, resp.Status)
	return nil
}

// Testa criação (upload) de um arquivo via acesso anônimo no file share.
func testFileUpload(account, share, fileName string) error {
	// URL para criar um arquivo no file share.
	url := fmt.Sprintf("https://%s.file.core.windows.net/%s/%s", account, share, fileName)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	// Cabeçalhos exigidos pela operação "Create File".
	// Exemplo: define o tamanho do arquivo que será criado.
	req.Header.Set("x-ms-content-length", "11") // tamanho em bytes para "upload teste"
	req.Header.Set("x-ms-type", "File")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("Upload anônimo do arquivo %q no file share %q retornou: %s\n", fileName, share, resp.Status)
	return nil
}

func main() {
	// Flags para definir o modo de teste e os parâmetros de conexão.
	mode := flag.String("mode", "blob", "Tipo de recurso para testar: 'blob' ou 'file'")
	account := flag.String("account", "", "Nome da conta de armazenamento")
	container := flag.String("container", "", "Nome do container (para teste de blob)")
	share := flag.String("share", "", "Nome do file share (para teste de file)")
	blobName := flag.String("blob", "testblob.txt", "Nome do blob para teste de upload")
	fileName := flag.String("file", "testfile.txt", "Nome do arquivo para teste de upload (file share)")
	flag.Parse()

	if *account == "" {
		fmt.Println("Informe o nome da conta usando a flag -account")
		os.Exit(1)
	}

	switch *mode {
	case "blob":
		if *container == "" {
			fmt.Println("Informe o nome do container usando a flag -container para teste de blob")
			os.Exit(1)
		}
		fmt.Println("Testando Azure Blob Storage via acesso anônimo:")
		if err := testBlobList(*account, *container); err != nil {
			fmt.Println("Erro ao listar container:", err)
		}
		if err := testBlobUpload(*account, *container, *blobName); err != nil {
			fmt.Println("Erro ao fazer upload do blob:", err)
		}
	case "file":
		if *share == "" {
			fmt.Println("Informe o nome do file share usando a flag -share para teste de file")
			os.Exit(1)
		}
		fmt.Println("Testando Azure File Shares via acesso anônimo:")
		if err := testFileList(*account, *share); err != nil {
			fmt.Println("Erro ao listar file share:", err)
		}
		if err := testFileUpload(*account, *share, *fileName); err != nil {
			fmt.Println("Erro ao fazer upload do arquivo:", err)
		}
	default:
		fmt.Println("Modo inválido. Utilize '-mode blob' ou '-mode file'")
	}
}
