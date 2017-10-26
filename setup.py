import setuptools
from packagename.version import Version


setuptools.setup(name='Open Redirect Scanner',
                 version=Version('1.0.0'),
                 description='Open Redirect Scanner',
                 author='Random_Robbie',
                 author_email='txt3rob@gmail.com',
                 url='http://path-to-my-packagename',
                 platforms='any',
                 install_requires=requirements,
                 license='MIT License',
                 zip_safe=False,
                 keywords='Open Redirect Scanner',
                 classifiers=['Packages', 'Open Redirect Scanner'])
