import unittest
import time
import random
import string
from unittest import mock
from shelloperator import executor


class ExecutorTestCase(unittest.TestCase):
    def setUp(self):
        self.run_mock = mock.MagicMock()

    def test_runs_function(self):
        with executor.Executor(concurrency=1, run_fn=self.run_mock) as x:
            x("sdsdff", 2)
            x("sdfggg", 4)
            x("54hdf", 6)

        time.sleep(0.3)
        self.assertEqual(self.run_mock.call_args_list, [mock.call(2), mock.call(4), mock.call(6)])


class HashingTestCase(unittest.TestCase):
    def test_hash_index_in_range(self):
        for _ in range(1000):
            to_hash = ''.join(random.choices(string.ascii_letters + string.digits, k=20))
            actual = executor.hash_to_index(to_hash, 3)

            self.assertLess(actual, 3)
            self.assertTrue(isinstance(actual, int))
