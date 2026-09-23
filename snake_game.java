import java.util.ArrayDeque;
import java.util.Deque;
import java.util.Objects;

class SnakeGame {
    int[][] board;
    int r;
    int c;
    Deque<int[]> snake;

    SnakeGame(int w, int h) {
        r = w;
        c = h;
        board = new int[r][c];
        snake = new ArrayDeque<>();
    }

    private boolean validateNewPosition(int i, int j) {
        if (i < 0 || j < 0 || i >= r || j >= c) return false;
        return board[i][j] != 2;
    }

    public void placeFood(int i, int j) {
        if (!validateNewPosition(i, j)) return;
        board[i][j] = 1;
    }

    private int[] getNewPosition(String command, int i, int j) {
        if (Objects.equals(command, "U")) {
            return new int[]{i - 1, j};
        } else if (Objects.equals(command, "R")) {
            return new int[]{i, j + 1};
        } else if (Objects.equals(command, "D")) {
            return new int[]{i + 1, j};
        } else if (Objects.equals(command, "L")) {
            return new int[]{i, j - 1};
        }
        return new int[]{i, j};
    }

    public int move(String command) {
        if (snake.isEmpty()) return -1;

        int[] currHead = snake.getFirst();
        int[] newHead = getNewPosition(command, currHead[0], currHead[1]);

        int[] tail = null;
        boolean isFood = (newHead[0] >= 0 && newHead[0] < r && newHead[1] >= 0 && newHead[1] < c && board[newHead[0]][newHead[1]] == 1);

        if (!isFood && !snake.isEmpty()) {
            tail = snake.pollLast();
            board[tail[0]][tail[1]] = 0;
        }

        if (!validateNewPosition(newHead[0], newHead[1])) {
            return -1;
        }

        if (isFood) {
            board[newHead[0]][newHead[1]] = 2;
            snake.addFirst(new int[]{newHead[0], newHead[1]});
        } else {
            board[newHead[0]][newHead[1]] = 2;
            snake.addFirst(new int[]{newHead[0], newHead[1]});
        }

        return snake.size();
    }
}
