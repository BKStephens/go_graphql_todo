import React from 'react';

const TodoList = ({todos, handleToggle}) => {
    const handleClick = (e) => {
        e.preventDefault();
        handleToggle(parseInt(e.currentTarget.id));
    }

    return (
        <div>
            {todos.map(todo => {
                return (
                    <div
                        style={todo.done ? styles.todoDone : styles.todo}
                        id={todo.id}
                        key={todo.id}
                        onClick={handleClick}
                    >
                        {todo.text}
                    </div>
                )
            })}
        </div>
    );
};

const styles = {
    todoDone: {
        cursor: 'pointer',
        textDecoration: 'line-through',
    },
    todo: {
        cursor: 'pointer',
    }
};
export default TodoList;
