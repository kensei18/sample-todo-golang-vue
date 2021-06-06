<template>
  <v-app>
    <v-container>
      <v-row>
        <v-col cols="4">
          <h1>ToDo</h1>
          <v-btn class="mb-2" color="primary" @click="showNewDialog"
            >Create</v-btn
          >
          <TaskList :task-list="todoList" @update-status="updateStatus(0)" />
        </v-col>

        <v-divider vertical></v-divider>

        <v-col cols="4">
          <h1>Doing</h1>
          <TaskList :task-list="doingList" @update-status="updateStatus(1)" />
        </v-col>

        <v-divider vertical></v-divider>

        <v-col cols="4">
          <h1>Done</h1>
          <TaskList :task-list="doneList" @update-status="updateStatus(2)" />
        </v-col>
      </v-row>
    </v-container>

    <v-dialog v-model="dialog" width="500">
      <v-card>
        <v-card-title>{{ dialogHeader }}</v-card-title>
        <v-card-text>
          <v-text-field v-model="name" label="Task Name"></v-text-field>
          <v-textarea
            v-model="description"
            label="Task Description"
            outlined
          ></v-textarea>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="green" dark @click="savingAction">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script>
import TaskList from "./TaskList";

export default {
  components: {
    TaskList,
  },

  provide() {
    return {
      updateSelectedTaskId: this.updateSelectedTaskId,
      showEditDialog: this.showEditDialog,
      removeTask: this.removeTask,
    };
  },

  data: () => ({
    tasks: [],
    dialog: false,
    selectedTaskId: 0,
    dialogHeader: "Create Task",
    savingAction: null,
    name: "",
    description: "",
  }),

  computed: {
    todoList() {
      return this.taskList(0);
    },
    doingList() {
      return this.taskList(1);
    },
    doneList() {
      return this.taskList(2);
    },
  },

  watch: {
    dialog(value) {
      if (!value) {
        this.name = "";
        this.description = "";
        this.selectedTaskId = 0;
      }
    },
  },

  mounted() {
    this.axios.get("/api/tasks").then(({ data }) => {
      data.forEach((task) => {
        task.status = Number(task.status);
      });
      this.tasks.push(...data);
    });
  },

  methods: {
    taskList(status) {
      return this.tasks.filter((task) => task.status === status);
    },
    findTaskById(id) {
      return this.tasks.find((task) => task.id === id);
    },
    updateSelectedTaskId(id) {
      this.selectedTaskId = id;
    },
    updateStatus(status) {
      const task = this.findTaskById(this.selectedTaskId);
      this.axios.put(`/api/tasks/${task.id}`, { status: status }).then(() => {
        task.status = status;
        this.selectedTaskId = 0;
      });
    },
    showNewDialog() {
      this.dialog = true;
      this.selectedTaskId = 0;
      this.dialogHeader = "Create Task";
      this.savingAction = this.createTask;
    },
    createTask() {
      this.axios
        .post("/api/tasks", {
          name: this.name,
          description: this.description,
        })
        .then(({ data }) => {
          this.tasks.push(data);
          this.dialog = false;
        });
    },
    showEditDialog(id) {
      const task = this.findTaskById(id);
      this.dialog = true;
      this.dialogHeader = "Edit Task";
      this.name = task.name;
      this.description = task.description;
      this.selectedTaskId = id;
      this.savingAction = this.updateTask;
    },
    updateTask() {
      const task = this.findTaskById(this.selectedTaskId);
      this.axios
        .put(`/api/tasks/${task.id}`, {
          name: this.name,
          description: this.description,
        })
        .then(() => {
          this.$set(task, "name", this.name);
          this.$set(task, "description", this.description);
          this.dialog = false;
        });
    },
    removeTask(id) {
      this.axios.delete(`/api/tasks/${id}`).then(() => {
        const idx = this.tasks.findIndex((task) => task.id === id);
        this.tasks.splice(idx, 1);
      });
    },
  },
};
</script>
