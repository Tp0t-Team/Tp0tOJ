<template>
  <v-container>
    <v-simple-table class="mb-4">
      <thead>
        <tr>
          <th class="text-left">Rank</th>
          <th class="text-left">Name</th>
          <th class="text-left">Score</th>
        </tr>
      </thead>
      <tbody>
        <tr
          class="table-item"
          v-for="r in pageValue"
          :key="r.index"
          @click="
            if ($store.state.global.role == 'admin')
              $router.push(`/profile/${r.desc.userId}`);
          "
        >
          <td>{{ r.index + 1 }}</td>
          <td>{{ r.desc.name }}</td>
          <td>{{ r.desc.score }}</td>
        </tr>
      </tbody>
    </v-simple-table>
    <v-row>
      <v-col cols="9" />
      <v-col cols="3">
        <v-text-field
          label="Search"
          outlined
          dense
          hide-details
          append-icon="search"
          @keydown.enter.prevent="updateFilter()"
          @click:append="updateFilter()"
          v-model="input"
        ></v-text-field>
      </v-col>
    </v-row>
    <v-row justify="center">
      <v-pagination
        v-model="page"
        :page="page"
        :length="pageCount"
      ></v-pagination>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { RankWithIndex } from "@/struct";

@Component
export default class RankTable extends Vue {
  @Prop({ required: true }) value!: RankWithIndex[];
  @Prop({ required: true }) userPerPage!: number;
  @Prop({ required: true }) pageInit!: number;

  private input: string = "";
  private filterRE: RegExp = /(.*)/;

  private updateFilter() {
    this.filterRE = new RegExp(
      ["", ...Array.from(this.input.matchAll(/./gu)).map(it => it[0]), ""].join(
        "(.*)"
      ),
      "iu"
    );
    this.page = 1;
  }

  private filter(it: RankWithIndex) {
    return this.filterRE.test(it.desc.name);
  }

  get displayValue(): RankWithIndex[] {
    return this.value.filter(it => this.filter(it));
  }

  get pageBase() {
    return (this.page - 1) * this.userPerPage;
  }

  get pageValue() {
    return this.displayValue.slice(
      this.pageBase,
      this.pageBase + this.userPerPage
    );
  }

  page: number = 1;
  get pageCount() {
    return Math.ceil(this.displayValue.length / this.userPerPage);
  }

  mounted() {
    this.page = this.pageInit;
  }
}
</script>

<style lang="scss" scoped></style>
