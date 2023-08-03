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

    def __init__(self):
        self.return_code = 0
        self.all_dockerfiles = []
        self.all_stages = []
        self.table = prettytable.PrettyTable()
        self.errors = []
        self.results = {}

    def add(self, top_element: str, dockerfile: str, stage: str, status: bool = True, message: str = None) -> None:
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
            self.errors.append('Failed {}/{} in {} stage, error message: {}'.format(
                top_element,
                dockerfile,
                stage,
                message,
            ))

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
            row = ['{}/{}'.format(*dockerfile)]
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
