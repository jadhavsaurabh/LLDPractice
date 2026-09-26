void main() {
    Solution s = new Solution();

    int[][] board = new int[][]{
            {0, 0, 1, 0, 2},
            {1, 1, 1, 2, 1},
            {2, 1, 1, 2, 2},
            {0, 0, 1, 0, 2},
            {1, 1, 1, 0, 0}
    };

    s.print(board);
    s.nextState(board);
}

class Solution {

    public void print(int[][] board) {
        System.out.println("__________________________");

        for (int[] row : board) {
            for (int value : row) {
                System.out.print(" " + value);
            }
            System.out.println();
        }
    }

    public void nextState(int[][] board) {

        while (true) {

            // Apply gravity
            applyGravity(board);
            print(board);

            // Mark pieces that should be removed
            boolean[][] marked = new boolean[board.length][board[0].length];

            boolean anythingMarked = mark(board, marked);

            // Nothing to remove, so the game is stable
            if (!anythingMarked) {
                break;
            }

            // Remove marked pieces
            sweep(board, marked);
            print(board);
        }
    }

    // Move all pieces down.
    private void applyGravity(int[][] board) {
        int r = board.length;
        int c = board[0].length;

        for (int j = 0; j < c; j++) {

            int ground = r - 1;

            for (int i = r - 1; i >= 0; i--) {

                if (board[i][j] != 0) {

                    if (ground != i) {
                        board[ground][j] = board[i][j];
                        board[i][j] = 0;
                    }

                    ground--;
                }
            }
        }
    }

    // Find all pieces that belong to a removable line.
    private boolean mark(int[][] board, boolean[][] marked) {

        int r = board.length;
        int c = board[0].length;

        boolean anythingMarked = false;

        for (int i = 0; i < r; i++) {
            for (int j = 0; j < c; j++) {

                if (board[i][j] == 0) {
                    continue;
                }

                int value = board[i][j];

                // Check vertical line.
                if (i > 0 && i < r - 1
                        && board[i - 1][j] == value
                        && board[i + 1][j] == value) {

                    marked[i][j] = true;

                    upWalk(board, marked, i, j);
                    downWalk(board, marked, i, j);

                    anythingMarked = true;
                }

                // Check horizontal line.
                if (j > 0 && j < c - 1
                        && board[i][j - 1] == value
                        && board[i][j + 1] == value) {

                    marked[i][j] = true;

                    leftWalk(board, marked, i, j);
                    rightWalk(board, marked, i, j);

                    anythingMarked = true;
                }
            }
        }

        return anythingMarked;
    }

    // Walk left and mark matching pieces.
    private void leftWalk(
            int[][] board,
            boolean[][] marked,
            int row,
            int col) {

        int value = board[row][col];

        col--;

        while (col >= 0 && board[row][col] == value) {
            marked[row][col] = true;
            col--;
        }
    }

    // Walk right and mark matching pieces.
    private void rightWalk(
            int[][] board,
            boolean[][] marked,
            int row,
            int col) {

        int value = board[row][col];

        col++;

        while (col < board[0].length && board[row][col] == value) {
            marked[row][col] = true;
            col++;
        }
    }

    // Walk up and mark matching pieces.
    private void upWalk(
            int[][] board,
            boolean[][] marked,
            int row,
            int col) {

        int value = board[row][col];

        row--;

        while (row >= 0 && board[row][col] == value) {
            marked[row][col] = true;
            row--;
        }
    }

    // Walk down and mark matching pieces.
    private void downWalk(
            int[][] board,
            boolean[][] marked,
            int row,
            int col) {

        int value = board[row][col];

        row++;

        while (row < board.length && board[row][col] == value) {
            marked[row][col] = true;
            row++;
        }
    }

    // Remove all marked pieces.
    private void sweep(int[][] board, boolean[][] marked) {

        for (int i = 0; i < board.length; i++) {
            for (int j = 0; j < board[0].length; j++) {

                if (marked[i][j]) {
                    board[i][j] = 0;
                }
            }
        }
    }
}
