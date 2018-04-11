import queue
import threading
import time
import hashlib


def hash_to_index(hash, ring):
    hshd = int(hashlib.md5(hash.encode('utf-8')).hexdigest(), 16)
    idx = hshd % ring
    return idx


def work(queue, run_fn):
    def do_work():
        while True:
            item = queue.get()

            if str(item) == "SHUTDOWN":
                return

            run_fn(*item)

    return do_work


class Executor(object):
    def __init__(self, *args, **kwargs):
        self.concurrency = kwargs['concurrency']
        self.run_fn = kwargs['run_fn']
        self.queues = []

    def __enter__(self):
        for i in range(self.concurrency):
            q = queue.Queue()
            t = threading.Thread(target=work(q, self.run_fn))

            self.queues.append((q, t))

            t.start()

        return self

    def __exit__(self, *args):
        for (q, t) in self.queues:
            if t.is_alive():
                q.put("SHUTDOWN")

        map(lambda _, t: t.join(), self.queues)

    def __call__(self, hash_key, *args):
        index = hash_to_index(hash_key, len(self.queues))
        (q, t) = self.queues[index]
        q.put(args)
