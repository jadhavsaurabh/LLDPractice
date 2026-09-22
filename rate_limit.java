void main() {
    RateLimiter r = new RateLimiter(3, 10);
    try {
        System.out.println("saurabh 1:" + r.isReqAllowed("saurabh"));
        System.out.println("saurabh 2:" + r.isReqAllowed("saurabh"));
        System.out.println("saurabh 3:" + r.isReqAllowed("saurabh"));
        System.out.println("saurabh 4:" + r.isReqAllowed("saurabh"));
        Thread.sleep(12000);

        System.out.println("saurabh 5:" + r.isReqAllowed("saurabh"));
    } catch (Exception d) {
        System.out.println("exception");
    }
}

class RateLimiter {
    int maxReqAllowed;
    int duration;
    HashMap<String, Queue<Long>> userHistory;
    private final Lock lock = new ReentrantLock();

    RateLimiter(int req, int d) {
        maxReqAllowed = req;
        duration = d;
        userHistory = new HashMap<>();
    }

    public boolean isReqAllowed(String userId) {
        long now = System.currentTimeMillis() / 1000;
        lock.lock();
        try {
            Queue<Long> userReq = userHistory.computeIfAbsent(userId, k -> new LinkedList<>());
            long windowStart = now - duration;

            while (!userReq.isEmpty() && userReq.peek() < windowStart) {
                userReq.poll();
            }

            if (userReq.size() < maxReqAllowed) {
                userReq.add(now);
                return true;
            }

            return false;
        } finally {
            lock.unlock();
        }
    }
}
