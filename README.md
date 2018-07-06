A simple AlphaGo-like AI for 5x5 [Hex](https://en.wikipedia.org/wiki/Hex_(board_game)), inspired by [Thinking Fast and Slow with Deep Learning and Tree Search](https://www.arxiv-vanity.com/papers/1705.08439/)

I'm building this to teach myself about machine learning. If it works well I'll experiment with the search algorithm, like [this repo](https://github.com/Videodr0me/leela-chess-experimental) has been doing for chess.

You can download training games and a trained model from [Google Drive](https://drive.google.com/drive/folders/1Q3zUscw3z6tTsOpR4WcFx1wbI11K8yLV?usp=sharing).

If you want to try it out yourself, here's how:

# Install dependencies

`dep ensure`

# Start Docker

Install [Docker Compose](https://docs.docker.com/compose/) and run:

```
$ docker-compose run --rm hexit
```

# Play an untrained AI

To play against an untrained AI, run `go run src/cmd/play/play.go`. The AI plays almost randomly, but if it can win in a single move it will.

```
- - - - -
 - - - - -
  - - - - -
   - - - - -
    - - - - -
0,0
X - - - -
 - - - - -
  - - - - -
   - - - - -
    - - O - -
1,1
X - - - -
 - X - - O
  - - - - -
   - - - - -
    - - O - -
```

# Generate training games

```
go run src/cmd/self_play/self_play.go
```

This will generate 1,000 training games in the `training_games/` folder.

# Train model

```
python3 neural/train.py
```

This will train a simple model and save it in the `hexit_saved_model/` folder.

# Test model against untrained AI

```
go run src/cmd/play_match/play_match.go
```

The trained model will play as Player 2, and on my machine it won 19 out of 20 games.
