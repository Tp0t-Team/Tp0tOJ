<template>
  <div>
    <v-container fluid fill-height class="scrollable">
      <v-row>
        <v-spacer></v-spacer>
        <v-col cols="8">
          <v-data-table :headers="headers" :items="images" :loading="loading">
            <template v-slot:item.digest="{ item }">
              <pre>{{ item.digest }}</pre>
            </template>
            <template v-slot:item.delete="{ item }">
              <v-btn
                text
                color="primary"
                @click="delImage(item.name)"
                :disabled="loading"
              >
                delete
              </v-btn>
            </template>
          </v-data-table>
          <!-- v-model="selected" -->
          <!-- item-key="name" -->
          <!-- show-select -->
        </v-col>
        <v-spacer></v-spacer>
      </v-row>
      <v-snackbar v-model="hasInfo" right bottom :timeout="3000">
        {{ infoText }}
        <!-- <v-spacer></v-spacer> -->
        <template v-slot:action="{ attrs }">
          <v-btn icon>
            <v-icon v-bind="attrs" @click="hasInfo = false">close</v-icon>
          </v-btn>
        </template>
      </v-snackbar>
    </v-container>
    <v-btn fab absolute right bottom color="primary" @click="enterEdit">
      <v-icon>add</v-icon>
    </v-btn>
    <v-dialog persistent v-model="edit" width="400px">
      <v-card width="400px">
        <v-form v-model="valid" ref="edit" class="ma-4">
          <v-row>
            <v-col cols="12">
              <v-file-input
                accept="application/x-tar"
                v-model="file"
                outlined
                hide-details
                prepend-icon=""
                prepend-inner-icon="attach_file"
                label="File Input"
              ></v-file-input>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="6">
              <v-text-field
                v-model="imageName"
                outlined
                hide-details
                label="Image Name"
              ></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="platform"
                outlined
                hide-details
                label="Platform"
              ></v-text-field>
            </v-col>
          </v-row>
        </v-form>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text :disabled="uploadLoading" @click="edit = false"
            >cancel</v-btn
          >
          <v-btn
            text
            :loading="uploadLoading"
            :disabled="uploadLoading"
            color="primary"
            @click="upload"
            >Upload</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import { ImageInfo, ImageInfoResult, Result } from "@/struct";

@Component
export default class Images extends Vue {
  private headers = [
    { text: "name", value: "name" },
    { text: "platform", value: "platform" },
    { text: "size", value: "size" },
    { text: "digest", value: "digest" },
    { text: "", value: "delete" }
  ];

  private selected: boolean[] = [];
  private images: ImageInfo[] = [];

  private valid: boolean = false;
  private edit: boolean = false;

  private loading: boolean = false;
  private uploadLoading: boolean = false;

  private imageName: string = "";
  private platform: string = "";
  private file: File | null = null;

  private infoText: string = "";
  private hasInfo: boolean = false;

  enterEdit() {
    this.edit = true;
  }

  async upload() {
    // this.edit = false;
    this.uploadLoading = true;
    try {
      const formData = new FormData();
      formData.set("name", this.imageName);
      formData.set("platform", this.platform);
      formData.append("image", this.file!);
      let res = await fetch("/image", {
        method: "POST",
        body: formData
      });
      if (res.status != 200) {
        throw res.statusText;
      }
      this.imageName = "";
      this.platform = "";
      this.file = null;
      await this.loadData();
      this.infoText = "success";
      this.hasInfo = true;
      await this.loadData();
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
    this.uploadLoading = false;
    this.edit = false;
    return true;
  }

  async mounted() {
    await this.loadData();
  }

  async loadData() {
    try {
      let res = await this.$apollo.query<ImageInfoResult, {}>({
        query: gql`
          query {
            imageInfos {
              message
              infos {
                name
                platform
                size
                digest
              }
            }
          }
        `,
        fetchPolicy: "no-cache"
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.imageInfos.message) throw res.data!.imageInfos.message;
      this.images = res.data!.imageInfos.infos.map(it => {
        it.digest = it.digest.slice(0, 8);
        console.log(it);
        let size = BigInt(it.size);
        if (size > 1024n * 1024n * 1024n) {
          it.size =
            (parseInt((size / 1024n / 1024n).toString()) / 1024).toFixed(2) +
            "GB";
        } else if (size > 1024n * 1024n) {
          it.size =
            (parseInt((size / 1024n).toString()) / 1024).toFixed(2) + "MB";
        } else if (size > 1024n) {
          it.size = (parseInt(it.size) / 1024).toFixed(2) + "KB";
        }
        return it;
      });
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  async delImage(name: string) {
    this.loading = true;
    try {
      let res = await this.$apollo.mutate<
        { deleteImage: Result },
        { input: string }
      >({
        mutation: gql`
          mutation($input: String!) {
            deleteImage(input: $input) {
              message
            }
          }
        `,
        variables: {
          input: name
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.deleteImage.message) throw res.data!.deleteImage.message;
      this.selected = [];
      this.loading = false;
      this.infoText = "delete success";
      this.hasInfo = true;
      await this.loadData();
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped>
.scrollable {
  overflow-y: auto;
}
</style>
