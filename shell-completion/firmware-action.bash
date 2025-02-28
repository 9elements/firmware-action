_firmware-action_spec_1 () {
    compgen -A file "$1"
}

_firmware-action () {
    if [[ $(type -t _get_comp_words_by_ref) != function ]]; then
        echo _get_comp_words_by_ref: function not defined.  Make sure the bash-completions system package is installed
        return 1
    fi

    local words cword
    _get_comp_words_by_ref -n "$COMP_WORDBREAKS" words cword

    local -a literals=("generate-config" "--config" "--json" "validate-config" "version" "build" "--indent" "--debug" "--target" "--recursive" "--interactive" "--help")

    declare -A literal_transitions
    literal_transitions[0]="([0]=1 [3]=1 [4]=2 [5]=3 [11]=2)"
    literal_transitions[1]="([1]=4 [11]=2)"
    literal_transitions[3]="([6]=5 [1]=6 [2]=5 [7]=5 [8]=7 [9]=5 [10]=5 [11]=2)"
    literal_transitions[5]="([6]=5 [1]=6 [2]=5 [7]=5 [8]=7 [9]=5 [10]=5)"

    declare -A match_anything_transitions
    match_anything_transitions=([4]=2 [7]=5 [6]=5)
    declare -A subword_transitions

    local state=0
    local word_index=1
    while [[ $word_index -lt $cword ]]; do
        local word=${words[$word_index]}

        if [[ -v "literal_transitions[$state]" ]]; then
            declare -A state_transitions
            eval "state_transitions=${literal_transitions[$state]}"

            local word_matched=0
            for literal_id in $(seq 0 $((${#literals[@]} - 1))); do
                if [[ ${literals[$literal_id]} = "$word" ]]; then
                    if [[ -v "state_transitions[$literal_id]" ]]; then
                        state=${state_transitions[$literal_id]}
                        word_index=$((word_index + 1))
                        word_matched=1
                        break
                    fi
                fi
            done
            if [[ $word_matched -ne 0 ]]; then
                continue
            fi
        fi

        if [[ -v "match_anything_transitions[$state]" ]]; then
            state=${match_anything_transitions[$state]}
            word_index=$((word_index + 1))
            continue
        fi

        return 1
    done


    local prefix="${words[$cword]}"

    local shortest_suffix="$word"
    for ((i=0; i < ${#COMP_WORDBREAKS}; i++)); do
        local char="${COMP_WORDBREAKS:$i:1}"
        local candidate="${word##*$char}"
        if [[ ${#candidate} -lt ${#shortest_suffix} ]]; then
            shortest_suffix=$candidate
        fi
    done
    local superfluous_prefix=""
    if [[ "$shortest_suffix" != "$word" ]]; then
        local superfluous_prefix=${word%$shortest_suffix}
    fi

    if [[ -v "literal_transitions[$state]" ]]; then
        local state_transitions_initializer=${literal_transitions[$state]}
        declare -A state_transitions
        eval "state_transitions=$state_transitions_initializer"

        for literal_id in "${!state_transitions[@]}"; do
            local literal="${literals[$literal_id]}"
            if [[ $literal = "${prefix}"* ]]; then
                local completion=${literal#"$superfluous_prefix"}
                COMPREPLY+=("$completion ")
            fi
        done
    fi

    declare -A specialized_commands
    specialized_commands=([4]=1 [6]=1)
    if [[ -v "specialized_commands[$state]" ]]; then
        local command_id=${specialized_commands[$state]}
        local completions=()
        mapfile -t completions < <(_firmware-action_spec_"${command_id}" "$prefix" | cut -f1)
        for item in "${completions[@]}"; do
            if [[ $item = "${prefix}"* ]]; then
                COMPREPLY+=("$item")
            fi
        done
    fi


    return 0
}

complete -o nospace -F _firmware-action firmware-action
