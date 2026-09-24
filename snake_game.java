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

        // Initial snake position
        snake.addFirst(new int[]{0, 0});
        board[0][0] = 2;
    }

    private boolean isValidPosition(int i, int j) {
        return i >= 0 && j >= 0 && i < r && j < c;
    }

    public void placeFood(int i, int j) {
        if (isValidPosition(i, j) && board[i][j] == 0) {
            board[i][j] = 1;
        }
    }

    private int[] getNewPosition(String command, int i, int j) {
        if (Objects.equals(command, "U")) {
            i--;
        } else if (Objects.equals(command, "R")) {
            j++;
        } else if (Objects.equals(command, "D")) {
            i++;
        } else if (Objects.equals(command, "L")) {
            j--;
        }

        return new int[]{i, j};
    }

    public int move(String command) {
        int[] currHead = snake.getFirst();
        int[] newHead = getNewPosition(command, currHead[0], currHead[1]);

        // Wall collision
        if (!isValidPosition(newHead[0], newHead[1])) {
            return -1;
        }

        boolean isFood = board[newHead[0]][newHead[1]] == 1;

        // Remove tail first if we are not eating.
        // This allows the head to move into the old tail position.
        if (!isFood) {
            int[] tail = snake.pollLast();
            board[tail[0]][tail[1]] = 0;
        }

        // Snake body collision
        if (board[newHead[0]][newHead[1]] == 2) {
            return -1;
        }

        // Add new head
        snake.addFirst(newHead);
        board[newHead[0]][newHead[1]] = 2;

        return snake.size();
    }
}
