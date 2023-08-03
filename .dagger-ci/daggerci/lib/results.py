'''
Class to handle results of building, testing and publishing
'''

import re
import logging
import prettytable


class Results():
    '''
    Class to handle results of building, testing and publishing
    '''

    def __init__(self) -> None:
        self.return_code = 0
        self.all_dockerfiles: list[list[str]] = []
        self.all_stages: list[str] = []
        self.table = prettytable.PrettyTable()
        self.errors: list[str] = []
        self.results: dict[str, dict[str, dict[str, str | bool | None]]] = {}

    def add(self,   # pylint: disable=too-many-arguments
            top_element: str,
            dockerfile: str,
            stage: str,
            status: bool = True,
            message: str | None = None) -> None:
        '''
        Add entry into results
        '''
        # Update list of all dockerfiles
        if [top_element, dockerfile] not in self.all_dockerfiles:
            self.all_dockerfiles.append([top_element, dockerfile])

        # Update list of all stages
        if stage not in self.all_stages and not re.match(r'.*_msg$', stage):
            self.all_stages.append(stage)

        # Update return code and collect error message
        if not status and message != 'skip':
            self.return_code = 1
            self.errors.append(
                f'Failed {top_element}/{dockerfile} in {stage} stage, error message: {message}')

        # Update results
        if top_element not in self.results:
            self.results[top_element] = {}
        if dockerfile not in self.results[top_element]:
            self.results[top_element][dockerfile] = {}
        self.results[top_element][dockerfile][stage] = status
        self.results[top_element][dockerfile][stage+'_msg'] = message

    def print(self) -> None:
        '''
        Pretty print the results
        '''
        # Construct the pretty table
        self.table.field_names = ['container']+self.all_stages
        for dockerfile in self.all_dockerfiles:
            row = ['{}/{}'.format(*dockerfile)]  # pylint: disable=consider-using-f-string
            for stage in self.all_stages:
                res = self.results[dockerfile[0]][dockerfile[1]]
                if stage in res:
                    if res[stage]:
                        row.append('OK')
                    else:
                        stmsg = stage+'_msg'
                        if stmsg in res and res[stmsg] == 'skip':
                            row.append('skip')
                        else:
                            row.append('fail')
                else:
                    row.append('-')
            self.table.add_row(row)

        # Print the prettytable
        print(self.table)

        # Print all found errors
        for i in self.errors:
            logging.error(i)
