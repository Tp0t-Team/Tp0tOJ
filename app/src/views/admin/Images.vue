<template>
  <div class="content-col">
    <v-container fill-width>
      <v-simple-table class="ma-4">
        <thead>
          <tr>
            <th class="text-left">Images</th>
            <th class="text-left">Size</th>
            <th class="text-left">Score</th>
            <th class="text-left">Operation</th>
          </tr>
        </thead>
        <tbody>
          <tr
            class="table-item"
            v-for="(r,index) in pageRank"
            :key="r.rank"
            @click="$router.push(`/profile/${r.userId}`)"
          >
            <td>{{ pageBase + index + 1 }}</td>
            <td>{{ r.name }}</td>
            <td>{{ r.score }}</td>
            <td>
              <v-btn>delete</v-btn>
            </td>
          </tr>
        </tbody>
      </v-simple-table>
      <v-row justify="center">
        <v-pagination v-model="page" :page="page" :length="pageCount"></v-pagination>
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
      <v-btn v-if="$store.state.global.role=='admin'||$store.state.global.role=='team'"
      fab
      absolute
      right
      bottom
      color="light-blue"
      @click="enterEdit"
    >
    <v-dialog v-model="edit" width="400px">
      <v-card width="400px" class="pa-4">
        <v-form v-model="valid" ref="edit">
          <v-row class="ma-4">
            <div>Dockerfile Upload</div>
          </v-row>
          <v-row class="pl-2">
            <v-text-field
              v-model="imageName"
              outlined
              label="ImageName"
            ></v-text-field>
          </v-row>
          <v-row class="pl-2">
            <v-file-input
              accept="application/x-tar"
              label="File input"
              v-model="file"
            ></v-file-input>
            <v-btn
              @click="upload"
            >Upload</v-btn>
          </v-row> 
        </v-form>
      </v-card>
    </v-dialog>
    <v-icon>add</v-icon>
  </v-btn>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import { RankDesc, RankResult } from "@/struct";

const UserPerPage = 10;

@Component
export default class Images extends Vue {
  private rankColor = ["amber", "light-blue", "green"];
  private page: number = 1;

  private ranks: RankDesc[] = [];
  private pageCount: number = 1;

  private infoText: string = "";
  private hasInfo: boolean = false;

  private valid: boolean = false;
  private loading: boolean = false;  
  private dockerfile:string = "helloworld";
  private edit:boolean = false;

  private imageName:string = "";
  private file:File|null =  null;

  enterEdit() {
    this.edit = true;
  }

  async upload(){
    this.edit = false;
    console.log(this.file);
    try {
      const formData = new FormData();
      formData.append('file', this.file!);
      let res = await fetch( '/image', {
        method: 'POST',
        headers: {
          'Content-Type': 'multipart/form-data'
        },
        body: formData
      })
      if (res.status != 200 ){
        throw res.statusText;
      }
      // let await res.text
    } catch (e) {
        console.log(e);
    }   
    return true;
  }

  private get topRank() {
    return this.ranks.slice(0, 3);
  }
  private get pageBase() {
    return (this.page - 1) * 10 + 3;
  }
  private get pageRank() {
    return this.ranks.slice(this.pageBase, this.pageBase + UserPerPage);
  }

  async mounted() {
    this.page = parseInt(this.$route.params.page);
    // try {
    //   let res = await this.$apollo.query<RankResult>({
    //     query: gql`
    //       query {
    //         rank {
    //           message
    //           rankResultDescs {
    //             userId
    //             name
    //             avatar
    //             score
    //           }
    //         }
    //       }
    //     `,
    //     fetchPolicy: "no-cache"
    //   });
    //   if (res.errors) throw res.errors.map(v => v.message).join(",");
    //   if (res.data!.rank.message) throw res.data!.rank.message;
    //   this.ranks = res.data!.rank.rankResultDescs.sort(
    //     (a, b) => parseInt(b.score) - parseInt(a.score)
    //   );
    //   this.pageCount = Math.floor(
    //     (this.ranks.length + UserPerPage - 1) / UserPerPage
    //   );
    // } catch (e) {
    //   this.infoText = e.toString();
    //   this.hasInfo = true;
    // }
  }
}
</script>

<style lang="scss" scoped>
.content-col {
  height: calc(100vh - 96px);
  overflow-y: auto;
}

.table-item {
  cursor: pointer;
}
</style>