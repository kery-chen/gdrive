package auth

import (
    "golang.org/x/oauth2"
    "encoding/json"
    "os"
    "io/ioutil"
)


func FileSource(path string, token *oauth2.Token, conf *oauth2.Config) oauth2.TokenSource {
    return &fileSource{
        tokenPath: path,
        tokenSource: conf.TokenSource(oauth2.NoContext, token),
    }
}

type fileSource struct {
    tokenPath string
    tokenSource oauth2.TokenSource
}

func (self *fileSource) Token() (*oauth2.Token, error) {
    token, err := self.tokenSource.Token()
    if err != nil {
        return token, err
    }

    // Save token to file
    SaveToken(self.tokenPath, token)

    return token, nil
}

func ReadToken(path string) (*oauth2.Token, bool, error) {
    if !fileExists(path) {
        return nil, false, nil
    }

    content, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, true, err
    }
    token := &oauth2.Token{}
    return token, true, json.Unmarshal(content, token)
}

func SaveToken(path string, token *oauth2.Token) error {
    data, err := json.MarshalIndent(token, "", "  ")
    if err != nil {
        return err
    }

    if err = mkdir(path); err != nil {
        return err
    }

    // Write to temp file first
    tmpFile := path + ".tmp"
    err = ioutil.WriteFile(tmpFile, data, 0600)
    if err != nil {
        os.Remove(tmpFile)
        return err
    }

    // Move file to correct path
    return os.Rename(tmpFile, path)
}
