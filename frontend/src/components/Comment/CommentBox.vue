<template>
	<div v-click-outside:[popperRef]="onClickOutside" class="form-box">
		<Editor
			ref="editorRef"
			v-model="content"
			:class="{ 'input-active': action }"
			:placeholder="props.placeholder"
			class="max-h-140px overflow-auto"
			:min-height="64"
			@focus="onFocus"
			@input="isEmpty(content.replace(/&nbsp; |<br>/g, '')) ? (disabled = true) : (disabled = false)"
			@submit="onSubmit"
		></Editor>
		<Transition name="fade">
			<div v-if="action && showBtn === true" class="action-box">
				<!-- <u-emoji :emoji="emoji" @add-emoji="(val: string) => editorRef?.addText(val)"></u-emoji> -->
				<el-button type="primary" :disabled="disabled" @click="onSubmit">
					{{ props.contentBtn }}
				</el-button>
			</div>
		</Transition>
	</div>
</template>

<script setup lang="ts">
import { isEmpty } from "../utils/index";
import { ClickOutside as vClickOutside } from "element-plus";
import { nextTick, ref } from "vue";
import Editor from "../Editor/index.vue";
// import { CommentSubmitParam, InjectionCommentFun } from "./interface";

export interface CommentBoxApi {
	focus(): void;
}

interface Props {
	placeholder: string;
	contentBtn: string;
	parentId?: number;
	replay?: string;
	showBtn?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
	placeholder: "请输入内容",
	contentBtn: "提交",
	showBtn: true
});
const content = ref("");
const action = ref(false);
const disabled = ref(true);
// const editorRef = ref<EditorInstance>();
const editorRef = ref();
const popperRef = ref();

const emit = defineEmits<{
	(e: "hide", event: Event): void;
}>();

// const submit = inject(InjectionCommentFun) as (obj: CommentSubmitParam) => void;
// const emoji = inject(InjectionEmojiApi);

const onSubmit = () => {
	let contentNew = props.replay ? `回复 <span style='color: blue;'>@${props.replay}:</span> ${content.value}` : content.value;
	console.log(contentNew);
	// submit({
	// 	clear: () => editorRef.value?.clear(),
	// 	content: props.replay ? `回复 <span style='color: blue;'>@${props.replay}:</span> ${content.value}` : content.value,
	// 	parentId: props.parentId
	// });
};

function onClickOutside(event: Event) {
	if (isEmpty(content.value)) {
		action.value = false;
		emit("hide", event);
	}
}

function onFocus() {
	action.value = true;
	nextTick(() => {
		popperRef.value = document.querySelector("div[id^='el-popper-container']");
	});
}

defineExpose({
	content,
	disabled,
	focus: () => (editorRef as any).value?.focus()
});
</script>

<style lang="scss" scoped>
.form-box {
	position: relative;
	width: 100%;
	padding: 1px;
	overflow: hidden;
	border-radius: 8px;
	.action-box {
		display: flex;
		align-items: center;
		margin-top: 8px;
		.el-button {
			margin-left: auto;
		}
	}
}
.fade-enter-active,
.fade-leave-active {
	transition: opacity 0.2s;
}
.fade-enter-from,
.fade-leave-to {
	opacity: 0;
}
</style>
