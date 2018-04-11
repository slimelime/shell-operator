import unittest
import tempfile
from unittest import mock
from shelloperator import commander


TEST_SCRIPT = """
#!/bin/bash

echo line 1

>&2 echo line stderr 1

sleep 2

echo line 2

echo line 3

>&2 echo line stderr 2

exit 44

"""

class CommanderTestCase(unittest.TestCase):
    def setUp(self):
        self.stdout_mock = mock.MagicMock()
        self.stderr_mock = mock.MagicMock()

    def test_run_command_get_stdout(self):
        exit_code = commander.run_command(["echo", "hello"], stdout_fn=self.stdout_mock, stderr_fn=self.stderr_mock)

        self.assertEqual(exit_code, 0)
        self.assertEqual(self.stderr_mock.call_args_list, [])
        self.assertEqual(self.stdout_mock.call_args_list, [mock.call(b"hello\n")])

    def test_run_script(self):
        with tempfile.NamedTemporaryFile() as t:
            t.write(TEST_SCRIPT.encode('utf-8'))
            t.flush()

            exit_code = commander.run_command(["/bin/bash", t.name], stdout_fn=self.stdout_mock, stderr_fn=self.stderr_mock)

            self.assertEqual(exit_code, 44)
            self.assertEqual(self.stderr_mock.call_args_list, [mock.call(b"line stderr 1\n"), mock.call(b"line stderr 2\n")])
            self.assertEqual(self.stdout_mock.call_args_list, [mock.call(b"line 1\n"), mock.call(b"line 2\n"), mock.call(b"line 3\n")])
