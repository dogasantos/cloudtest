package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Testa a listagem de objetos de um bucket via acesso anônimo.
func testGCSList(bucket string) error {
	// URL da API JSON para listar objetos:
	// https://storage.googleapis.com/storage/v1/b/{bucket}/o
	url := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/%s/o", bucket)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Listagem anônima do bucket %q retornou: %s\n", bucket, resp.Status)
	fmt.Printf("Resposta: %s\n", string(body))
	return nil
}

// Testa a leitura de um objeto específico no bucket via acesso anônimo.
func testGCSRead(bucket, object string) error {
	// URL para acessar diretamente um objeto:
	// https://storage.googleapis.com/{bucket}/{object}
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, object)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Leitura anônima do objeto %q no bucket %q retornou: %s\n", object, bucket, resp.Status)
	fmt.Printf("Conteúdo: %s\n", string(body))
	return nil
}

// Testa o upload de um objeto para o bucket via acesso anônimo.
func testGCSUpload(bucket, object string) error {
	// URL para upload simples de objeto:
	// https://storage.googleapis.com/upload/storage/v1/b/{bucket}/o?uploadType=media&name={object}
	url := fmt.Sprintf("https://storage.googleapis.com/upload/storage/v1/b/%s/o?uploadType=media&name=%s", bucket, object)
	data := "conteúdo de teste"
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Upload anônimo do objeto %q no bucket %q retornou: %s\n", object, bucket, resp.Status)
	fmt.Printf("Resposta: %s\n", string(body))
	return nil
}

func main() {
	mode := flag.String("mode", "list", "Tipo de teste: 'list', 'read' ou 'upload'")
	bucket := flag.String("bucket", "", "Nome do bucket do GCS")
	object := flag.String("object", "test.txt", "Nome do objeto para teste de leitura ou upload")
	flag.Parse()

	if *bucket == "" {
		fmt.Println("Informe o nome do bucket usando a flag -bucket")
		os.Exit(1)
	}

	switch *mode {
	case "list":
		fmt.Println("Testando listagem anônima do bucket...")
		if err := testGCSList(*bucket); err != nil {
			fmt.Println("Erro ao listar bucket:", err)
		}
	case "read":
		fmt.Println("Testando leitura anônima do objeto...")
		if err := testGCSRead(*bucket, *object); err != nil {
			fmt.Println("Erro ao ler objeto:", err)
		}
	case "upload":
		fmt.Println("Testando upload anônimo para o bucket...")
		if err := testGCSUpload(*bucket, *object); err != nil {
			fmt.Println("Erro ao fazer upload do objeto:", err)
		}
	default:
		fmt.Println("Modo inválido. Utilize 'list', 'read' ou 'upload'")
	}
}
