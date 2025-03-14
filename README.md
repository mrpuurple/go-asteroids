# Go Asteroids

## Gameplay

- start the game (spacebar)
- fire (spacebar)
- move forward (keyUp)
- move backwards (keyDown)
- activate shield (keyS)
- quit the game in the "Game Over" scene (keyQ)

## Init

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

### Run the program

```sh
go run .
```

**If you encounter**<br>
`current directory is contained in a module that is not one of the workspace modules listed in go.work.`<br>
You can add the module to the workspace using:

```sh
go work use .
```
