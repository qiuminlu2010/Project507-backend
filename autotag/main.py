import torch
from src.models import create_model
import numpy as np
from PIL import Image
from src.helper_functions.helper_functions import parse_args

class Args:
    def __init__(self,model_path,model_name,dataset_type,th=None,input_size=448):
        self.model_path = model_path
        self.model_name = model_name
        self.th = th
        self.dataset_type = dataset_type
        self.input_size = input_size

def __init__():
    args = Args(
        model_name='tresnet_l',
        model_path = "pretrained/Open_ImagesV6_TRresNet_L_448.pth", 
        dataset_type = 'OpenImages',
        input_size= 448
    )
    args = parse_args(args)
    # 加载模型
    state = torch.load(args.model_path, map_location='cpu')
    args.num_classes = state['num_classes']
    model = create_model(args).cuda()
    model.load_state_dict(state['model'], strict=True)
    model.eval()
    labels = []
    with open("label.txt","r") as f:
        for d in f.readlines():
            labels.append(d.strip())

    classes_list = np.array(labels)

    # 读取图片
    pic_path = "../微信图片_20220804012910.jpg"
    im = Image.open(pic_path)
    im_resize = im.resize((args.input_size, args.input_size))
    np_img = np.array(im_resize, dtype=np.uint8)
    tensor_img = torch.from_numpy(np_img).permute(2, 0, 1).float() / 255.0  # HWC to CHW
    tensor_batch = torch.unsqueeze(tensor_img, 0).cuda()

    # 推理标签
    output = torch.squeeze(torch.sigmoid(model(tensor_batch)))
    np_output = output.cpu().detach().numpy()
    detected_classes = classes_list[np_output > args.th]

    