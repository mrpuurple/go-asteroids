# Go Asteroids

- [Go Asteroids](#go-asteroids)
  - [Gameplay](#gameplay)
  - [Initial setup](#initial-setup)
    - [Get dependencies](#get-dependencies)
      - [Deploy section (optional)](#deploy-section-optional)
    - [Create the code](#create-the-code)
    - [Start coding 😊](#start-coding-)
    - [Run the program](#run-the-program)
    - [Package](#package)
      - [Create a Taskfile (needs to go in your root folder)](#create-a-taskfile-needs-to-go-in-your-root-folder)
      - [Test it](#test-it)
    - [Note](#note)

## Gameplay

| Action | Key |
|--------|-----|
| Start the game | spacebar |
| Fire | spacebar |
| Move forward | keyUp |
| Move backwards | keyDown |
| Rotate left | keyLeft |
| Rotate Right | keyRight |
| Activate shield | keyS |
| Activate HyperSpace | keyH |
| Quit the game in the "Game Over" scene | keyQ |

## Initial setup

```sh
go mod init asteroids
```

### Get dependencies

```sh
go get github.com/hajimehoshi/ebiten/v2
go get github.com/solarlune/resolv@v0.8
```

#### Deploy section (optional)

```sh
go install github.com/machinebox/appify@latest
```

### Create the code

```sh
touch main.go
```

### Start coding 😊

### Run the program

```sh
go run .
```

### Package

```sh
go install github.com/machinebox/appify@latest
```

```sh
which appify
/Users/ssscse/go/bin/appify
```

#### Create a Taskfile (needs to go in your root folder)

```yaml
version: '3'


tasks:

  run:
    aliases:
      - start
      - go
      - run
    cmds:
      - go run .

  build:
    cmds:
      - go build -o dist/asteroids .

  dmg:
    cmds:
      - task appify
      - rm -rf dist/Go\ Asteroids.app/
      - mv "Go Asteroids.app" dist
      - rm -f dist/asteroids
      - hdiutil create -volname "Go Asteroids" -srcfolder dist -ov -format UDZO Asteroids.dmg

  build-wasm:
    cmds:
      - env GOOS=js GOARCH=wasm go build -o wasm/asteroids.wasm .

  serve:
    deps: [build-wasm]
    cmds:
      - cd wasm && python -m http.server 4000

  appify:
    - task: build
    - appify -author "Nichlas Handstedt" -version "v1.0.0" -name "Go Asteroids" -icon ./assets/images/icon.png ./dist/asteroids

  clean:
    cmds:
      - go clean
      - rm -rf dist/*
      - rm -f Asteroids.dmg
      - rm -rf "Go Asteroids.app"
```

#### Test it

```sh
task run
```

### Note

If you encounter

`current directory is contained in a module that is not one of the workspace modules listed in go.work.`<br>

You can add the module to the workspace using:

```sh
go work use .
```
