# Copyright 2020-Present Mark Spicer
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
# documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
# rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
# permit persons to whom the Software is furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
# Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
# WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
# OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

load("cache.star", "cache")
load("encoding/base64.star", "base64")
load("render.star", "render")
load("time.star", "time")


"""
Config.
"""
DEV = False                 # Set to true to skip the intro and go straight to game play.
VICTORY = 1000              # The number of generations to win.
START_LIVE_CELLS = 400      # The number of starting live cells to populate your society with.


"""
Constants.
"""
ALIVE = 1
DEAD = 0
PIXEL_WIDTH = 64
PIXEL_HEIGHT = 32
FRAMES_PER_VIEW = 60
BOARD_PADDING = 10
BOARD_WIDTH = PIXEL_WIDTH + BOARD_PADDING
BOARD_HEIGHT = PIXEL_HEIGHT + BOARD_PADDING


"""
Cache keys.
"""
KEY_STATE = "state"
KEY_BOARD = "board"
KEY_GEN = "generation"


"""
State machine.
"""
STATE_STARTING = "starting"
STATE_RUNNING = "running"
STATE_GAME_OVER = "game_over"
STATE_VICTORY = "victory"
STATES = [
    STATE_STARTING,
    STATE_RUNNING,
    STATE_GAME_OVER,
    STATE_VICTORY,
]


"""
Game runtime.
"""
def main():
    """
    Welcome to the Tidbyt rendition of Conway's Game of Life! This app will render generations of life in between your
    other Tidbyt apps. It starts with a start of game display and then will render generations every time the app cycles
    on until either your society has survived past the defined generations or your society has died out! In which case,
    either a vicory message or game over message will appear during the next app cycle respectively.
    """
    if DEV:
        # Skip the introduction if we are in development.
        setup_game()

    state = get_state()
    if state == STATE_STARTING:
        return start_game()

    if state == STATE_VICTORY or state == STATE_GAME_OVER:
        return end_game(state)

    return run_game()


def start_game():
    """
    Starts a new game by loading the first board into the cache and rendering the start message.
    """
    setup_game()
    return render_start()


def run_game():
    """
    Run the game as normal. If there is a victory condition, render the remaining frames anyways to play out the game.
    If there is a game over condition due to loops or an empty board, set game over and render a skull for the rest of
    the frames.
    """
    board = get_board()
    if board == None:
        return start_game()

    gen = get_generation()
    state = get_state()
    empty_board = create_board(BOARD_WIDTH, BOARD_HEIGHT)
    frames = []
    boards = [encode(board)]

    loops = False
    for x in range(FRAMES_PER_VIEW):
        # If there is an empty board, we can short circuit additional rendering.
        if state == STATE_GAME_OVER and not loops:
            frames.append(render_frame(board_subset(board)))
            continue
        
        # Add our frame, generate our next board.
        frames.append(render_frame(board_subset(board)))
        board = next_generation(board)

        # If our society has died and we have not survived past the victory condition, set game over.
        if board == empty_board and state != STATE_VICTORY:
            state = set_state(STATE_GAME_OVER)
            continue
            
        # If our board exists in the last 30 frames, we can assume there is a loop. These are called oscillators, and
        # they can exist at of periods 4, 8, 14, 15, 30. Given we are generating 30 frames, we should check all 30 for
        # an oscillating society.
        encoded = encode(board)
        if encoded in boards and state != STATE_VICTORY:
            loops = True
            state = set_state(STATE_GAME_OVER)
            continue
        
        # Add out board to the list.
        boards.append(encoded)

        # Set victory state if we have survived as a society long enough! Generate the rest of the frames anyways to see
        # extended cuts. 
        gen += 1
        if gen > VICTORY:
            state = set_state(STATE_VICTORY)
            continue

    # Add our last board generation to the cache to snag it next time.
    set_board(board)
    set_generation(gen)

    return render.Root(
        delay = 500,
        child = render.Animation(
            children=frames,
        ),
    )


def end_game(state):
    """
    Ends the game by setting the state back to starting and rendering the victory or game over message according to
    the supplied state.
    """
    set_state(STATE_STARTING)
    gen = get_generation()

    if state == STATE_GAME_OVER:
        return render_game_over(gen)

    return render_victory(gen)


def setup_game():
    """
    Sets up a new game be rendering a board and setting initial cache values.
    """
    board = create_board(BOARD_WIDTH, BOARD_HEIGHT)
    seed_board(board, BOARD_WIDTH, BOARD_HEIGHT, START_LIVE_CELLS)
    set_board(board)
    set_generation(0)
    set_state(STATE_RUNNING)


"""
Cache getters/setters.
"""
def get_state():
    """
    Gets the state out of the cache. Returns the starting state if it does not exist.
    """
    cached_state = cache.get(KEY_STATE)

    if cached_state in STATES:
        return cached_state

    return STATE_STARTING


def set_state(state):
    """
    Sets the state in the cache and returns it so you can track it in the calling function.
    """
    cache.set(KEY_STATE, state, ttl_seconds=600)
    return state


def get_generation():
    """
    Getter for the generation out of the cache. Returns 0 if it does not exist.
    """
    cached_generation = cache.get(KEY_GEN)
    if cached_generation == None:
        return 0

    return int(cached_generation)


def set_generation(gen):
    """
    Setter for the generation in the cache.
    """
    cache.set(KEY_GEN, str(gen), ttl_seconds=300)


def get_board():
    """
    Getter for the board in the cache. Returns none if it doesn't exist.
    """
    cached_board = cache.get(KEY_BOARD)
    if cached_board == None:
        return None

    return decode(cached_board)


def set_board(board):
    """
    Setter for the board in the cache.
    """
    cache.set(KEY_BOARD, encode(board), ttl_seconds=300)


def encode(board):
    """
    Encodes the board into a string to be stored in the cache.
    """
    return "\n".join([
        ",".join([str(cell) for cell in row])
        for row in board
    ])


def decode(cached_board):
    """
    Decodes the board from the cache back into a 2D array from a string.
    """
    return [
        [int(x) for x in line.split(",")]
        for line in cached_board.split("\n")
    ]


"""
Game mechanics.
"""
def next_generation(board):
    """
    Creates the next generation based off of the existing board.
    """
    return [
        [dead_or_alive(board, x, y) for y, cell in enumerate(row)]
        for x, row in enumerate(board)
    ]


def board_subset(board):
    """
    Returns a subset of the board to render into frames. This is useful so we can avoid edges in our board.
    """
    pad = int(BOARD_PADDING / 2)
    return [
        [cell for cell in row[pad:BOARD_WIDTH-pad]]
        for row in board[pad:BOARD_HEIGHT-pad]
    ]


def dead_or_alive(board, x, y):
    """
    Determines if a cell should be dead or alive, based off of the rules found here:
    https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life

    - Any live cell with two or three live neighbours survives.
    - Any dead cell with three live neighbours becomes a live cell.
    - All other live cells die in the next generation. Similarly, all other dead cells stay dead.
    """
    num_neighbors = calc_num_neighbors(board, x, y)

    # Survives.
    if board[x][y] == ALIVE and (num_neighbors == 2 or num_neighbors == 3):
        return ALIVE

    # Reproduce.
    if board[x][y] == DEAD and num_neighbors == 3:
        return ALIVE
    
    # Dies or was already dead.
    return DEAD


def calc_num_neighbors(board, x, y):
    """
    Returns the number of neighbors given a board and the position to calculate.
    """
    return len([
        1
        for i in range(max(0, x-1), min(x+2, len(board)))
        for j in range(max(0, y-1), min(y+2, len(board[i])))
        if x != i or y != j
        if board[i][j] == ALIVE
    ])


def create_board(width, height):
    """
    Creates a board using the provided witdth and height.
    """
    return [
        [DEAD for y in range(width)] 
        for x in range(height)
    ] 


def seed_board(board, width, height, num_start):
    """
    Modifies a given board with the desired number of starting live cells, placed randomly.
    """
    current = 0
    seed = time.now().nanosecond

    while current < num_start:
        x, seed = random(seed)
        y, seed = random(seed)

        x = int(x * 1000) % height
        y = int(y * 1000) % width

        if board[x][y] == DEAD:
            board[x][y] = ALIVE
            current = current + 1


"""
Rendering.
"""
def render_victory(gen):
    """
    Renders victory message for the game.
    """
    return render.Root(
        child=render.Box(
            child=render.Column(
                main_align="space_around",
                cross_align="center",
                children=[
                    render.Text(
                        content="Victory!",
                        font="5x8",
                    ),
                    render.Marquee(
                        child=render.Text(
                            content="Survived Generations: " + str(gen),
                            font="5x8",
                        ),
                        width=64,
                    ),
                ],
            ),
        )
    )


def render_game_over(gen):
    """
    Renders game over message for the game.
    """
    return render.Root(
        child=render.Box(
            child=render.Column(
                main_align="space_around",
                cross_align="center",
                children=[
                    render.Text(
                        content="Game Over.",
                        font="5x8",
                    ),
                    render.Marquee(
                        child=render.Text(
                            content="Survived Generations: " + str(gen),
                            font="5x8",
                        ),
                        width=64,
                    ),
                ],
            ),
        )
    )


def render_start():
    """
    Renders the start message for the game.
    """
    board = get_seeded_gosper_glider_board()

    frames = []
    for x in range(FRAMES_PER_VIEW):
        frames.append(render_frame(board_subset(board), text="Start"))
        board = next_generation(board)

    return render.Root(
        delay = 500,
        child = render.Animation(
            children=frames,
        ),
    )


def render_frame(board, text=None):
    """
    Renders a frame for a given board. Use this in an animation to display each round.
    """
    children = [
        render.Column(
            children=[render_row(row) for row in board],
        ),
    ]

    if text:
        children.append(
            render.Padding(
                child=render.Text(
                    content=text,
                    font="6x13",
                ),
                pad=(2,19,0,0),
            )
        )
    return render.Stack(
        children=children,
    )


def render_row(row):
    """
    Helper to render a row.
    """
    return render.Row(children=[render_cell(cell) for cell in row])


def render_cell(cell):
    """
    Helper to render a cell.
    """
    color = "#aaa" if cell == ALIVE else "#000"
    return render.Box(width=1, height=1, color=color)


"""
Utilities.
"""
def random(seed):
    """
    Returns a random number and the new seed value.
    
    Starlark is meant to be deterministic, so anything that made the language non-deterministic (such as random number
    generators) was removed. This is a Python implementation of a linear congruential generator I found here:
    http://www.cs.wm.edu/~va/software/park/park.html
    """
    modulus = 2147483648
    multiplier = 48271

    q = modulus / multiplier
    r = modulus % multiplier
    t = multiplier * (seed % q) - r * (seed / q);

    if t > 0:
        seed = t
    else:
        seed = t + modulus

    return float(seed / modulus), seed


def get_seeded_gosper_glider_board():
    """
    Returns a board rendered with an active Gosper Glider Gun. The board is encoded to a string using our encode()
    method and then base64 encoded. This was a trade off in repeatability with efficiency. I originally rendered this
    by using points mapped to create the Gosper Glider Gun and running through a few iterations until it looked nice.
    """
    encoded_board = "MCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwxLDEsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwxLDAsMCwxLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDEsMCwxLDAsMCwwLDAsMCwwLDAsMCwwLDEsMCwwLDAsMCwwLDAsMCwxLDEsMSwxLDEsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwxLDAsMCwwLDEsMCwwLDAsMCwwLDAsMCwxLDAsMCwwLDAsMCwwLDEsMCwwLDAsMCwwLDEsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwxLDAsMCwwLDAsMCwwLDAsMSwwLDAsMCwwLDAsMCwwLDEsMSwwLDAsMCwxLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwxLDEsMCwwLDAsMCwwLDAsMCwwLDEsMCwwLDAsMCwxLDAsMCwwLDAsMCwwLDAsMSwwLDAsMSwwLDAsMCwwLDAsMCwwLDEsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMSwxLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDEsMCwwLDAsMCwwLDAsMCwwLDAsMCwxLDEsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMSwwLDAsMCwxLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDEsMCwxLDAsMCwwLDAsMCwwLDAsMCwxLDAsMSwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwxLDEsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMSwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDEsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDEsMSwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMSwxLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwxLDAsMSwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwxLDEsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMSwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAKMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwCjAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMAowLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDA="
    return decode(base64.decode(encoded_board))