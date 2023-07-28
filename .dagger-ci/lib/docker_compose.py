'''
Functions to help parse Docker Compose
Docs: https://docs.docker.com/compose/compose-file/

!!! This is class / parses is incomplete implementation of Docker Compose specification !!!
It implements only specific functions that are needed in this top_element.
'''

import yaml
import dagger


def select(heap: list, needle: str = None) -> str:
    '''
    Select either top_element or dockerfile from conpose YAML
        - as default return first defined top_element or dockerfile
        - return needle
    '''
    # Default
    if needle is None:
        return heap[0]
    # If needle exists in heap, return it
    if needle in heap:
        return needle
    # If needle does not exists, raise exception
    raise ValueError


class DockerCompose():
    '''
    Class to parse docker compose file
    '''

    def __init__(self, path: str):
        self.path = path
        with open(path, 'r', encoding='utf-8') as composefile:
            self.yaml = yaml.safe_load(composefile.read())
            # yaml.safe_load might raise yaml.YAMLError

    def get_top_elements(self) -> list:
        '''
        Return a list of all top_elements found (list of strings)
        '''
        return list(self.yaml.keys())

    def __select_top_element__(self, top_element: str = None) -> str:
        '''
        Check if top_element in YAML
          if true, return said top_element name
          if None provided, default to the first top_element in YAML
          if false, raise ValueError exception
        '''
        try:
            return select(heap=self.get_top_elements(), needle=top_element)
        except ValueError:
            raise ValueError(f'Top element {top_element} not found in YAML file')

    def get_dockerfiles(self, top_element: str = None) -> list:
        '''
        Return a list of all docker files in top_element (list of strings)
        if no top_element provided, use the first one
        '''
        this_top_element = self.__select_top_element__(top_element)
        return list(self.yaml[this_top_element].keys())

    def __select_dockerfile__(self, dockerfile: str = None, top_element: str = None) -> str:
        '''
        Check if dockerfile under top_element in YAML
          if true, return said dockerfile name
          if None provided, default to the first dockerfile under top_element
          if false, raise ValueError exception
        '''
        this_top_element = self.__select_top_element__(top_element)
        try:
            return select(heap=self.get_dockerfiles(top_element=this_top_element), needle=dockerfile)
        except ValueError:
            raise ValueError(
                f'Dockerfile {dockerfile} not found in YAML file under {top_element} top_element')

    def get_dockerfile_context(self, dockerfile: str = None, top_element: str = None) -> str:
        '''
        Return a context of given dockerfile
        '''
        this_top_element = self.__select_top_element__(top_element)
        this_dockerfile = self.__select_dockerfile__(dockerfile)

        if 'build' in self.yaml[this_top_element][this_dockerfile]:
            if 'context' in self.yaml[this_top_element][this_dockerfile]['build']:
                return self.yaml[this_top_element][this_dockerfile]['build']['context']
        return None

    def get_dockerfile_args(self, dockerfile: str = None, top_element: str = None) -> list:
        '''
        Return a list of args for given dockerfile
            return list of dagger.api.gen.BuildArg
            https://dagger-io.readthedocs.io/en/sdk-python-v0.6.4/api.html#dagger.api.gen.BuildArg
        '''
        this_top_element = self.__select_top_element__(top_element)
        this_dockerfile = self.__select_dockerfile__(dockerfile)

        if 'args' in self.yaml[this_top_element][this_dockerfile]['build']:
            return [dagger.api.gen.BuildArg(i.split('=')[0], i.split('=')[1]) for i in self.yaml[this_top_element][this_dockerfile]['build']['args']]
        return []
