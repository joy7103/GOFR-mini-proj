import React from "react";
import Tasks from "./Tasks";
import { Paper, TextField } from "@material-ui/core";
import { Checkbox, Button } from "@material-ui/core";
import "./App.css";

class App extends Tasks {
    state = { tasks: [], currentTask: "" };
    render() {
        const { tasks } = this.state;
        return (
            <div className="App flex">
                <Paper elevation={3} className="container">
                    <div className="heading">TO-DO</div>
                    <form
                        onSubmit={this.handleSubmit}
                        className="flex"
                        style={{ margin: "15px 0" }}
                    >
                        <TextField
                            variant="outlined"
                            size="small"
                            style={{ width: "80%" }}
                            value={this.state.currentTask}
                            required={true}
                            onChange={this.handleChange}
                            placeholder="Add New TO-DO"
                        />
                        <Button
                            style={{ height: "40px" }}
                            color="primary"
                            variant="outlined"
                            type="submit"
                        >
                            Add task
                        </Button>
                    </form>
                    <div>
                    {tasks && tasks.length > 0 ? (
    tasks.map((task) => (
        <Paper key={task.id} className="flex task_container">
            <Checkbox
                checked={task.completed}
                onClick={() => this.handleUpdate(task.id)}
                color="primary"
            />
            <div
                className={
                    task.completed
                        ? "task line_through"
                        : "task"
                }
            >
                {task.task}
            </div>
            <Button
                onClick={() => {
                    console.log("Deleting task with ID:", task.id);
                    this.handleDelete(task.id);
                }}
                color="secondary"
            >
                delete
            </Button>
        </Paper>
    ))
) : (
    <div>No tasks found. Add a new task above.</div>
)}
                    </div>
                </Paper>
            </div>
        );
    }
}

export default App;
