/** @jsx h */

<div>
    <h3>TODO</h3>
    <TodoList items={state.items} />
    <form onSubmit={handleSubmit}>
        <label htmlFor="new-todo">
            What needs to be done?
          </label>
        <input
            id="new-todo"
            onChange={handleChange}
            value={state.text}
        />
        <button>
            Add #{state.items.length + 1}
        </button>
    </form>
</div>
