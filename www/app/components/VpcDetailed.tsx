/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as VpcTypes from '../types/VpcTypes';
import * as VpcActions from '../actions/VpcActions';
import * as OrganizationTypes from "../types/OrganizationTypes";
import * as PageInfos from './PageInfo';
import DatacentersStore from "../stores/DatacentersStore";
import OrganizationsStore from "../stores/OrganizationsStore";
import VpcRoute from './VpcRoute';
import VpcSubnet from './VpcSubnet';
import PageInput from './PageInput';
import PageInfo from './PageInfo';
import PageSave from './PageSave';
import ConfirmButton from './ConfirmButton';
import Help from './Help';
import PageTextArea from "./PageTextArea";

interface Props {
	organizations: OrganizationTypes.OrganizationsRo;
	vpc: VpcTypes.VpcRo;
	selected: boolean;
	onSelect: (shift: boolean) => void;
	onClose: () => void;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	addNetworkRole: string;
	addVpc: string;
	vpc: VpcTypes.Vpc;
}

const css = {
	card: {
		position: 'relative',
		padding: '48px 10px 0 10px',
		width: '100%',
	} as React.CSSProperties,
	button: {
		height: '30px',
	} as React.CSSProperties,
	buttons: {
		cursor: 'pointer',
		position: 'absolute',
		top: 0,
		left: 0,
		right: 0,
		padding: '4px',
		height: '39px',
		backgroundColor: 'rgba(0, 0, 0, 0.13)',
	} as React.CSSProperties,
	item: {
		margin: '9px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
	itemsLabel: {
		display: 'block',
	} as React.CSSProperties,
	itemsAdd: {
		margin: '8px 0 15px 0',
	} as React.CSSProperties,
	group: {
		flex: 1,
		minWidth: '280px',
		margin: '0 10px',
	} as React.CSSProperties,
	save: {
		paddingBottom: '10px',
	} as React.CSSProperties,
	label: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	status: {
		margin: '6px 0 0 1px',
	} as React.CSSProperties,
	icon: {
		marginRight: '3px',
	} as React.CSSProperties,
	inputGroup: {
		width: '100%',
	} as React.CSSProperties,
	protocol: {
		flex: '0 1 auto',
	} as React.CSSProperties,
	port: {
		flex: '1',
	} as React.CSSProperties,
	select: {
		margin: '7px 0px 0px 6px',
		paddingTop: '3px',
	} as React.CSSProperties,
	role: {
		margin: '9px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
	list: {
		marginBottom: '15px',
	} as React.CSSProperties,
};

export default class VpcDetailed extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			addNetworkRole: null,
			addVpc: null,
			vpc: null,
		};
	}

	set(name: string, val: any): void {
		let vpc: any;

		if (this.state.changed) {
			vpc = {
				...this.state.vpc,
			};
		} else {
			vpc = {
				...this.props.vpc,
			};
		}

		vpc[name] = val;

		this.setState({
			...this.state,
			changed: true,
			vpc: vpc,
		});
	}

	onAddSubnet = (i: number): void => {
		let vpc: VpcTypes.Vpc;

		if (this.state.changed) {
			vpc = {
				...this.state.vpc,
			};
		} else {
			vpc = {
				...this.props.vpc,
			};
		}

		let subnets = [
			...(vpc.subnets || []),
		];

		if (subnets.length === 0) {
			subnets = [{}];
		}

		subnets.splice(i + 1, 0, {} as VpcTypes.Subnet);
		vpc.subnets = subnets;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			vpc: vpc,
		});
	}

	onChangeSubnet(i: number, subnet: VpcTypes.Subnet): void {
		let vpc: VpcTypes.Vpc;

		if (this.state.changed) {
			vpc = {
				...this.state.vpc,
			};
		} else {
			vpc = {
				...this.props.vpc,
			};
		}

		let subnets = [
			...(vpc.subnets || []),
		];

		if (subnets.length === 0) {
			subnets = [{}];
		}

		subnets[i] = subnet;

		vpc.subnets = subnets;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			vpc: vpc,
		});
	}

	onRemoveSubnet(i: number): void {
		let vpc: VpcTypes.Vpc;

		if (this.state.changed) {
			vpc = {
				...this.state.vpc,
			};
		} else {
			vpc = {
				...this.props.vpc,
			};
		}

		let subnets = [
			...(vpc.subnets || []),
		];

		if (subnets.length !== 0) {
			subnets.splice(i, 1);
		}

		if (subnets.length === 0) {
			subnets = [{}];
		}

		vpc.subnets = subnets;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			vpc: vpc,
		});
	}

	onAddRoute = (i: number): void => {
		let vpc: VpcTypes.Vpc;

		if (this.state.changed) {
			vpc = {
				...this.state.vpc,
			};
		} else {
			vpc = {
				...this.props.vpc,
			};
		}

		let routes = [
			...(vpc.routes || []),
		];

		routes.splice(i + 1, 0, {} as VpcTypes.Route);
		vpc.routes = routes;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			vpc: vpc,
		});
	}

	onChangeRoute(i: number, route: VpcTypes.Route): void {
		let vpc: VpcTypes.Vpc;

		if (this.state.changed) {
			vpc = {
				...this.state.vpc,
			};
		} else {
			vpc = {
				...this.props.vpc,
			};
		}

		let routes = [
			...vpc.routes,
		];

		routes[i] = route;

		vpc.routes = routes;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			vpc: vpc,
		});
	}

	onRemoveRoute(i: number): void {
		let vpc: VpcTypes.Vpc;

		if (this.state.changed) {
			vpc = {
				...this.state.vpc,
			};
		} else {
			vpc = {
				...this.props.vpc,
			};
		}

		let routes = [
			...vpc.routes,
		];

		routes.splice(i, 1);

		vpc.routes = routes;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			vpc: vpc,
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		VpcActions.commit(this.state.vpc).then((): void => {
			this.setState({
				...this.state,
				message: 'Your changes have been saved',
				changed: false,
				disabled: false,
			});

			setTimeout((): void => {
				if (!this.state.changed) {
					this.setState({
						...this.state,
						vpc: null,
						changed: false,
					});
				}
			}, 1000);

			setTimeout((): void => {
				if (!this.state.changed) {
					this.setState({
						...this.state,
						message: '',
					});
				}
			}, 3000);
		}).catch((): void => {
			this.setState({
				...this.state,
				message: '',
				disabled: false,
			});
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		VpcActions.remove(this.props.vpc.id).then((): void => {
			this.setState({
				...this.state,
				disabled: false,
			});
		}).catch((): void => {
			this.setState({
				...this.state,
				disabled: false,
			});
		});
	}

	render(): JSX.Element {
		let vpc: VpcTypes.Vpc = this.state.vpc ||
			this.props.vpc;

		let datacenter = DatacentersStore.datacenter(vpc.datacenter);
		let datacenterName = datacenter ? datacenter.name : vpc.datacenter;
		let org = OrganizationsStore.organization(this.props.vpc.organization);

		let subnets = (vpc.subnets || []);
		if (subnets.length === 0) {
			subnets.push({});
		}

		let subnetsElem: JSX.Element[] = [];
		for (let i = 0; i < subnets.length; i++) {
			let index = i;

			subnetsElem.push(
				<VpcSubnet
					key={index}
					subnet={subnets[index]}
					onChange={(state: VpcTypes.Subnet): void => {
						this.onChangeSubnet(index, state);
					}}
					onAdd={(): void => {
						this.onAddSubnet(index);
					}}
					onRemove={(): void => {
						this.onRemoveSubnet(index);
					}}
				/>,
			);
		}

		let routes: JSX.Element[] = [
			<VpcRoute
				disabled={true}
				key={-1}
				route={{
					destination: '0.0.0.0/0',
					target: '0.0.0.0',
				} as VpcTypes.Route}
				onAdd={(): void => {
					this.onAddRoute(-1);
				}}
			/>,
		];
		if (vpc.routes) {
			for (let i = 0; i < (vpc.routes || []).length; i++) {
				let index = i;

				routes.push(
					<VpcRoute
						key={index}
						route={vpc.routes[index]}
						onChange={(state: VpcTypes.Route): void => {
							this.onChangeRoute(index, state);
						}}
						onAdd={(): void => {
							this.onAddRoute(index);
						}}
						onRemove={(): void => {
							this.onRemoveRoute(index);
						}}
					/>,
				);
			}
		}

		let fields: PageInfos.Field[] = [
			{
				label: 'ID',
				value: this.props.vpc.id || 'Unknown',
			},
			{
				label: 'Datacenter',
				value: datacenterName,
			},
			{
				label: 'Organization',
				value: org ? org.name : this.props.vpc.organization,
			},
			{
				label: 'VLAN Number',
				value: this.props.vpc.vpc_id || 'Unknown',
			},
			{
				label: 'Private IPv6 Network',
				value: this.props.vpc.network6 || 'Unknown',
				copy: true,
			},
		];

		return <td
			className="bp3-cell"
			colSpan={6}
			style={css.card}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div
						className="layout horizontal"
						style={css.buttons}
						onClick={(evt): void => {
							let target = evt.target as HTMLElement;

							if (target.className.indexOf('open-ignore') !== -1) {
								return;
							}

							this.props.onClose();
						}}
					>
            <div>
              <label
                className="bp3-control bp3-checkbox open-ignore"
                style={css.select}
              >
                <input
                  type="checkbox"
                  className="open-ignore"
                  checked={this.props.selected}
									onChange={(evt): void => {
									}}
                  onClick={(evt): void => {
										this.props.onSelect(evt.shiftKey);
									}}
                />
                <span className="bp3-control-indicator open-ignore"/>
              </label>
            </div>
						<div className="flex"/>
						<ConfirmButton
							className="bp3-minimal bp3-intent-danger bp3-icon-trash open-ignore"
							style={css.button}
							safe={true}
							progressClassName="bp3-intent-danger"
							dialogClassName="bp3-intent-danger bp3-icon-delete"
							dialogLabel="Delete VPC"
							confirmMsg="Permanently delete this VPC"
							confirmInput={true}
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						label="Name"
						help="Name of vpc"
						type="text"
						placeholder="Enter name"
						value={vpc.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<PageTextArea
						label="Comment"
						help="VPC comment."
						placeholder="VPC comment"
						rows={3}
						value={vpc.comment}
						onChange={(val: string): void => {
							this.set('comment', val);
						}}
					/>
					<PageInput
						label="Network"
						help="Network address of vpc with cidr."
						type="text"
						placeholder="Enter network"
						value={vpc.network}
						onChange={(val): void => {
							this.set('network', val);
						}}
					/>
					<label style={css.itemsLabel}>
						Subnets
						<Help
							title="Subnets"
							content="Subnets in VPC, can only be added or removed. Once added a subnet network block cannot be modified."
						/>
					</label>
					<div style={css.list}>
						{subnetsElem}
					</div>
				</div>
				<div style={css.group}>
					<PageInfo
						fields={fields}
					/>
					<label style={css.itemsLabel}>
						Route Table
						<Help
							title="Route Table"
							content="VPC routing table, enter a CIDR network for the desitnation and IP address for taget."
						/>
					</label>
					<div style={css.list}>
						{routes}
					</div>
				</div>
			</div>
			<PageSave
				style={css.save}
				hidden={!this.state.vpc && !this.state.message}
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						vpc: null,
					});
				}}
				onSave={this.onSave}
			/>
		</td>;
	}
}
