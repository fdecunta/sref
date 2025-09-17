import unittest

from refprintf import abbreviate_names


class TestFormatMethods(unittest.TestCase):
    
    def test_abbreviate_names(self):
        names = "John A."
        short_names = abbreviate_names(names)
        self.assertEqual(short_names, "J. A.")

        names = "Homero J."
        short_names = abbreviate_names(names)
        self.assertEqual(short_names, "H. J.")
