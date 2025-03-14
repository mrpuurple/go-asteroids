# Go Asteroids

- [Go Asteroids](#go-asteroids)
  - [Gameplay](#gameplay)
  - [Initial setup](#initial-setup)
    - [Get dependencies](#get-dependencies)
    - [Create the program](#create-the-program)
    - [Start coding 😊](#start-coding-)
    - [Run the program](#run-the-program)
    - [Note](#note)

## Gameplay

| Action | Key |
|--------|-----|
| Start the game | spacebar |
| Fire | spacebar |
| Move forward | keyUp |
| Move backwards | keyDown |
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

### Create the program

```sh
touch main.go
```

### Start coding 😊

### Run the program

```sh
go run .
```

### Note

If you encounter

`current directory is contained in a module that is not one of the workspace modules listed in go.work.`<br>

You can add the module to the workspace using:

```sh
go work use .
```
